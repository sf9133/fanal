package applier

import (
	"sort"
	"testing"

	godeptypes "github.com/sf9133/go-dep-parser/pkg/types"
	"github.com/stretchr/testify/assert"

	"github.com/sf9133/fanal/types"
)

func TestApplyLayers(t *testing.T) {
	testCases := []struct {
		name                   string
		inputLayers            []types.BlobInfo
		expectedArtifactDetail types.ArtifactDetail
	}{
		{
			name: "happy path",
			inputLayers: []types.BlobInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					OS: &types.OS{
						Family: "alpine",
						Name:   "3.10",
					},
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
							},
						},
					},
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "gemlibrary1",
										Version: "1.2.3",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
								{ // added
									Name:    "musl",
									Version: "1.2.4",
									Release: "4.5.7",
								},
							},
						},
					},
					WhiteoutFiles: []string{"app/composer.lock"},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "lib/apk/db/installed",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
								{
									Name:    "musl",
									Version: "1.2.4",
									Release: "4.5.8", // updated
								},
							},
						},
					},
				},
			},
			expectedArtifactDetail: types.ArtifactDetail{
				OS: &types.OS{
					Family: "alpine",
					Name:   "3.10",
				},
				Packages: []types.Package{
					{
						Name:    "musl",
						Version: "1.2.4",
						Release: "4.5.8",
						Layer: types.Layer{
							Digest: "sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
					{
						Name:    "openssl",
						Version: "1.2.3",
						Release: "4.5.6",
						Layer: types.Layer{
							Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
				},
				Applications: []types.Application{
					{
						Type:     "gem",
						FilePath: "app/Gemfile.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "gemlibrary1",
									Version: "1.2.3",
								},
								Layer: types.Layer{
									Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
									DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "happy path with removed and updated lockfile",
			inputLayers: []types.BlobInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					OS: &types.OS{
						Family: "alpine",
						Name:   "3.10",
					},
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "rails",
										Version: "5.0.0",
									},
								},
								{
									Library: godeptypes.Library{
										Name:    "rack",
										Version: "4.0.0",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					Applications: []types.Application{
						{
							Type:     "gem",
							FilePath: "app/Gemfile.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "rails",
										Version: "6.0.0",
									},
								},
								{
									Library: godeptypes.Library{
										Name:    "rack",
										Version: "4.0.0",
									},
								},
							},
						},
						{
							Type:     "composer",
							FilePath: "app/composer2.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
					WhiteoutFiles: []string{"app/composer.lock"},
				},
			},
			expectedArtifactDetail: types.ArtifactDetail{
				OS: &types.OS{
					Family: "alpine",
					Name:   "3.10",
				},
				Applications: []types.Application{
					{
						Type:     "gem",
						FilePath: "app/Gemfile.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "rack",
									Version: "4.0.0",
								},
								Layer: types.Layer{
									Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
									DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
								},
							},
							{
								Library: godeptypes.Library{
									Name:    "rails",
									Version: "6.0.0",
								},
								Layer: types.Layer{
									Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
									DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
								},
							},
						},
					},
					{
						Type:     "composer",
						FilePath: "app/composer2.lock",
						Libraries: []types.LibraryInfo{
							{
								Library: godeptypes.Library{
									Name:    "phplibrary1",
									Version: "6.6.6",
								},
								Layer: types.Layer{
									Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
									DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "happy path with status.d",
			inputLayers: []types.BlobInfo{
				{
					SchemaVersion: 1,
					Digest:        "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
					DiffID:        "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
					OS: &types.OS{
						Family: "debian",
						Name:   "8",
					},
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "var/lib/dpkg/status.d/openssl",
							Packages: []types.Package{
								{
									Name:    "openssl",
									Version: "1.2.3",
									Release: "4.5.6",
								},
							},
						},
					},
					Applications: []types.Application{
						{
							Type:     "composer",
							FilePath: "app/composer.lock",
							Libraries: []types.LibraryInfo{
								{
									Library: godeptypes.Library{
										Name:    "phplibrary1",
										Version: "6.6.6",
									},
								},
							},
						},
					},
				},
				{
					SchemaVersion: 1,
					Digest:        "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
					DiffID:        "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
					PackageInfos: []types.PackageInfo{
						{
							FilePath: "var/lib/dpkg/status.d/libc",
							Packages: []types.Package{
								{
									Name:    "libc",
									Version: "1.2.4",
									Release: "4.5.7",
								},
							},
						},
					},
					OpaqueDirs: []string{"app"},
				},
			},
			expectedArtifactDetail: types.ArtifactDetail{
				OS: &types.OS{
					Family: "debian",
					Name:   "8",
				},
				Packages: []types.Package{
					{
						Name:    "libc",
						Version: "1.2.4",
						Release: "4.5.7",
						Layer: types.Layer{
							Digest: "sha256:932da51564135c98a49a34a193d6cd363d8fa4184d957fde16c9d8527b3f3b02",
							DiffID: "sha256:a187dde48cd289ac374ad8539930628314bc581a481cdb41409c9289419ddb72",
						},
					},
					{
						Name:    "openssl",
						Version: "1.2.3",
						Release: "4.5.6",
						Layer: types.Layer{
							Digest: "sha256:24df0d4e20c0f42d3703bf1f1db2bdd77346c7956f74f423603d651e8e5ae8a7",
							DiffID: "sha256:aad63a9339440e7c3e1fff2b988991b9bfb81280042fa7f39a5e327023056819",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotArtifactDetail := ApplyLayers(tc.inputLayers)
			sort.Slice(gotArtifactDetail.Packages, func(i, j int) bool {
				return gotArtifactDetail.Packages[i].Name < gotArtifactDetail.Packages[j].Name
			})
			sort.Slice(gotArtifactDetail.Applications, func(i, j int) bool {
				return gotArtifactDetail.Applications[i].FilePath < gotArtifactDetail.Applications[j].FilePath
			})
			for _, app := range gotArtifactDetail.Applications {
				sort.Slice(app.Libraries, func(i, j int) bool {
					return app.Libraries[i].Library.Name < app.Libraries[j].Library.Name
				})
			}
			assert.Equal(t, tc.expectedArtifactDetail, gotArtifactDetail, tc.name)
		})
	}
}
