package lib

import (
	"context"
	crypto_rand "crypto/rand"
	_ "embed"
	"encoding/base64"
	"os"
	"path"
	"strconv"

	"github.com/gitrules/gitrules/github/common"
	"github.com/gitrules/gitrules/gitrules/api"
	"github.com/gitrules/gitrules/lib/base"
	"github.com/gitrules/gitrules/lib/form"
	"github.com/gitrules/gitrules/lib/git"
	"github.com/gitrules/gitrules/lib/must"
	"github.com/gitrules/gitrules/lib/ns"
	gitprovider "github.com/gitrules/gitrules/lib/provider"
	"github.com/gitrules/gitrules/proto/boot"
	"github.com/gitrules/gitrules/proto/gov"
	"github.com/google/go-github/v66/github"
	"golang.org/x/crypto/nacl/box"
)

func Deploy(
	ctx context.Context,
	token string, // permissions: read project issues, create/write govPrefix
	project common.Repo,
	govPrefix common.Repo,
	release string, // GitHub release of GitRules to install
) api.Config {

	ghc := common.NewClientForToken(ctx, token)

	pp, idOwnerAddr := common.CreatePublicPrivateRepos(ctx, ghc, project, govPrefix)
	govOwnerAddr := gov.OwnerAddress(idOwnerAddr)

	// attach access token authentication to context for git use
	git.SetAuth(ctx, govOwnerAddr.Public.Repo, git.MakeTokenAuth(ctx, token))
	git.SetAuth(ctx, govOwnerAddr.Private.Repo, git.MakeTokenAuth(ctx, token))

	// initialize governance identity
	base.Infof("initializing organizational governance for %v", project)
	boot.Boot(ctx, govOwnerAddr)

	// create GitHub environment for governance
	base.Infof("creating GitHub environment for governance in %v", pp.Public)
	createDeployEnvironment(ctx, ghc, token, project, pp.Public, pp.PublicURLs, pp.PrivateURLs, release)

	// install github automation in the public governance repo
	base.Infof("installing GitHub actions for governance in %v, targetting %v", pp.Public, project)
	installGithubActions(ctx, govOwnerAddr)

	// install governance labels in project repo
	createGovernanceIssueLabels(ctx, ghc, project)

	// return config for gitrules administrator
	homeDir, err := os.UserHomeDir()
	must.NoError(ctx, err)
	return api.Config{
		Auth: map[git.URL]api.AuthConfig{
			git.URL(pp.PublicURLs.HTTPSURL):               {AccessToken: github.String(token)},
			git.URL(pp.PrivateURLs.HTTPSURL):              {AccessToken: github.String(token)},
			git.URL("YOUR_MEMBER_PUBLIC_REPO_HTTPS_URL"):  {AccessToken: github.String("YOUR_MEMBER_ACCESS_TOKEN")},
			git.URL("YOUR_MEMBER_PRIVATE_REPO_HTTPS_URL"): {AccessToken: github.String("YOUR_MEMBER_ACCESS_TOKEN")},
		},
		//
		GovPublicURL:     git.URL(pp.PublicURLs.HTTPSURL),
		GovPublicBranch:  git.MainBranch,
		GovPrivateURL:    git.URL(pp.PrivateURLs.HTTPSURL),
		GovPrivateBranch: git.MainBranch,
		//
		MemberPublicURL:     "YOUR_MEMBER_PUBLIC_REPO_HTTPS_URL",
		MemberPublicBranch:  git.MainBranch,
		MemberPrivateURL:    "YOUR_MEMBER_PRIVATE_REPO_HTTPS_URL",
		MemberPrivateBranch: git.MainBranch,
		//
		CacheDir:        path.Join(homeDir, ".gitrules", "cache"),
		CacheTTLSeconds: 0,
	}
}

var (
	//go:embed deploy/.github/scripts/gitrules_cron.sh
	cronSH string

	//go:embed deploy/.github/workflows/gitrules_cron.yml
	cronYML string

	//go:embed deploy/.github/python/requirements.txt
	pythonRequirementsTXT string
)

func installGithubActions(
	ctx context.Context,
	govOwnerAddr gov.OwnerAddress,
) {

	govCloned := git.CloneOne(ctx, git.Address(govOwnerAddr.Public))
	t := govCloned.Tree()

	// populate helper files for github actions
	git.StringToFileStage(ctx, t, ns.NS{".github", "scripts", "gitrules_cron.sh"}, cronSH)
	git.StringToFileStage(ctx, t, ns.NS{".github", "workflows", "gitrules_cron.yml"}, cronYML)
	git.StringToFileStage(ctx, t, ns.NS{".github", "python", "requirements.txt"}, pythonRequirementsTXT)

	git.Commit(ctx, t, "install gitrules github actions")
	govCloned.Push(ctx)
}

func createGovernanceIssueLabels(
	ctx context.Context,
	ghc *github.Client,
	project common.Repo,
) {

	for _, l := range GovernanceLabels {
		label := &github.Label{Name: github.String(l)}
		_, _, err := ghc.Issues.CreateLabel(ctx, project.Owner, project.Name, label)
		if IsLabelAlreadyExists(err) {
			base.Infof("github issue label %v already exists in %v", l, project)
			continue
		}
		must.NoError(ctx, err)
	}
}

func createDeployEnvironment(
	ctx context.Context,
	ghClient *github.Client,
	token string,
	project common.Repo,
	govPublic common.Repo,
	govPublicURLs *gitprovider.Repository,
	govPrivateURLs *gitprovider.Repository,
	ghRelease string,
) {

	// fetch repo id
	ghGovPubRepo, _, err := ghClient.Repositories.Get(ctx, govPublic.Owner, govPublic.Name)
	must.NoError(ctx, err)

	// create deploy environment
	createEnv := &github.CreateUpdateEnvironment{}
	env, _, err := ghClient.Repositories.CreateUpdateEnvironment(ctx, govPublic.Owner, govPublic.Name, DeployEnvName, createEnv)
	must.NoError(ctx, err)

	// create environment secrets
	envSecrets := map[string]string{
		DeployEnvAccessToken: token,
	}

	govEnvPubKey, _, err := ghClient.Actions.GetEnvPublicKey(ctx, int(ghGovPubRepo.GetID()), env.GetName())
	// govPubPubKey, _, err := ghClient.Actions.GetRepoPublicKey(ctx, govPublic.Owner, govPublic.Name)
	must.NoError(ctx, err)

	for k, v := range envSecrets {
		encryptedValue := EncryptSecret(ctx, govEnvPubKey, v)
		encryptedSecret := &github.EncryptedSecret{
			Name:           k,
			KeyID:          govEnvPubKey.GetKeyID(),
			EncryptedValue: encryptedValue,
		}
		base.Infof("adding secret to environment: %v", form.SprintJSON(encryptedSecret))
		_, err := ghClient.Actions.CreateOrUpdateEnvSecret(ctx, int(ghGovPubRepo.GetID()), env.GetName(), encryptedSecret)
		must.NoError(ctx, err)
	}

	// create environment variables
	envVars := map[string]string{
		"GITRULES_RELEASE":    ghRelease,
		"PROJECT_OWNER":       project.Owner,
		"PROJECT_REPO":        project.Name,
		"PUBLIC_REPO_URL":     govPublicURLs.HTTPSURL,
		"PRIVATE_REPO_URL":    govPrivateURLs.HTTPSURL,
		"SYNC_GITHUB_FREQ":    strconv.Itoa(DefaultGithubFreq),
		"SYNC_COMMUNITY_FREQ": strconv.Itoa(DefaultCommunityFreq),
		"SYNC_FETCH_PAR":      strconv.Itoa(DefaultFetchParallelism),
	}
	for k, v := range envVars {
		_, err := ghClient.Actions.CreateEnvVariable(ctx, govPublic.Owner, govPublic.Name, env.GetName(), &github.ActionsVariable{Name: k, Value: v})
		must.NoError(ctx, err)
	}
}

const (
	DefaultGithubFreq       = 120 // seconds
	DefaultCommunityFreq    = 120 // seconds
	DefaultFetchParallelism = 5
)

func EncryptSecret(ctx context.Context, pubKey *github.PublicKey, secretValue string) string {

	decodedPubKey, err := base64.StdEncoding.DecodeString(pubKey.GetKey())
	must.NoError(ctx, err)

	var boxKey [32]byte
	copy(boxKey[:], decodedPubKey)
	secretBytes := []byte(secretValue)
	encryptedBytes, err := box.SealAnonymous([]byte{}, secretBytes, &boxKey, crypto_rand.Reader)
	must.NoError(ctx, err)

	return base64.StdEncoding.EncodeToString(encryptedBytes)
}
