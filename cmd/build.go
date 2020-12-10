/*
Copyright © 2020 DSCVIT

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/GDGVIT/katamari/internal/utils"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build your katamari project",
	Long:  `Fetch all repos from the specified organization, clone the READMEs and generate static pages ready for hosting!`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		org := viper.GetString("site")

		isUser := viper.GetBool("isuser")

		if org == "" {
			utils.Err("config", fmt.Sprintf("configuration not found! make sure you have a %s file in the "+
				"project dir.", chalk.Green.Color(".katamari.toml")))
			os.Exit(1)
		}

		utils.Info("sill", fmt.Sprintf("fetching repos for %s", chalk.Green.Color(fmt.Sprintf("'%s'", org))))

		err := os.MkdirAll("./content/readmes/", 0755)
		if err != nil {
			utils.Err("enoent", err.Error())
		}

		ctx := context.Background()

		token := os.Getenv("GITHUB_ACCESS_TOKEN")

		var client *github.Client

		if token != "" {
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)
			client = github.NewClient(tc)
		} else {
			utils.Warn("config", "Access Token Missing. Set an access token in a .katamari/config.json, else you might be rate limited by GitHub [refer README]")
			client = github.NewClient(nil)
		}

		orgOpt := &github.RepositoryListByOrgOptions{
			Type:        "public",
			ListOptions: github.ListOptions{PerPage: 50},
		}

		usrOpt := &github.RepositoryListOptions{
			Type:        "public",
			ListOptions: github.ListOptions{PerPage: 50},
		}

		var allRepos []*github.Repository

		if isUser {
			for {
				repos, resp, err := client.Repositories.List(ctx, org, usrOpt)
				if err != nil {
					utils.Err("enoent", err.Error())
				}
				allRepos = append(allRepos, repos...)
				if resp.NextPage == 0 {
					break
				}
				usrOpt.Page = resp.NextPage
			}
		} else {
			for {
				repos, resp, err := client.Repositories.ListByOrg(ctx, org, orgOpt)
				if err != nil {
					utils.Err("enoent", err.Error())
				}
				allRepos = append(allRepos, repos...)
				if resp.NextPage == 0 {
					break
				}
				orgOpt.Page = resp.NextPage
			}
		}

		if err != nil {
			utils.Err("enoent", err.Error())
			os.Exit(1)
		}

		var wg sync.WaitGroup

		for _, repo := range allRepos {
			wg.Add(1)
			utils.Info("sill", fmt.Sprintf("Fetching readme for repo %s", *repo.Name))
			go func(client *github.Client, repo *github.Repository, wg *sync.WaitGroup) {
				defer wg.Done()
				readme, _, err := client.Repositories.GetReadme(context.Background(), org, *repo.Name, nil)
				if err != nil {
					utils.Err("enoent", err.Error())
					return
				}

				content, err := readme.GetContent()
				if err != nil {
					utils.Err("enoent", err.Error())
					return
				}

				f, err := os.Create(fmt.Sprintf("./content/readmes/%s.md", *repo.Name))
				if err != nil {
					utils.Err("enoent", err.Error())
					return
				}

				_, _ = f.WriteString("---\n")
				_, _ = f.WriteString(fmt.Sprintf("title: %s\n", *repo.Name))
				_, _ = f.WriteString(fmt.Sprintf("date: %s\n", time.Now().UTC().Format("2006-01-02T15:04:05-0700")))
				_, _ = f.WriteString(fmt.Sprintf("draft: false\n"))
				_, _ = f.WriteString("---\n")
				_, _ = f.WriteString(content)
			}(client, repo, &wg)
		}
		wg.Wait()

		utils.Info("sill", fmt.Sprintf("Fetched %d repos", len(allRepos)))
		utils.Info("sill", chalk.Green.Color("Successfully built your katamari project!"))
		utils.Info("sill", fmt.Sprintf("Run %s %s", chalk.Green.Color("hugo server"), chalk.White.Color("to run the hugo server")))
		utils.Info("sill", fmt.Sprintf("To change theme refer %s", chalk.Green.Color("hugo documentation")))
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
