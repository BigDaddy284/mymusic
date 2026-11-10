package main

import "github.com/egoist/mygo/ui"

// The groups of Settings › About: this copy of MyTube and its licenses.

func init() { addSettings("about", 10, (*shell).mytubeSettings) }

// mytubeSettings tell of this copy of MyTube, with a button that opens its
// logs.
func (s *shell) mytubeSettings(c *ui.Context) {
	info := s.svc.System.Info()
	s.group(c, "MyTube", func() {
		version := "MyTube " + info.Version + " · " + info.Platform + " " + info.Arch
		if build, ok := buildNames[s.svc.System.InstallKind()]; ok {
			version += " · " + build
		}
		s.row(c, "Version", version, func() {
			if s.button(c, "logs", "open logs folder", false).Clicked() {
				if err := s.svc.System.OpenLogsFolder(); err != nil {
					s.notes = append(s.notes, noticeOf(err))
				}
			}
		}, "logs folder")
	})
}

func init() { addSettings("about", 30, (*shell).licenseSettings) }

// licenseSettings show the licenses of what MyTube is made of.
func (s *shell) licenseSettings(c *ui.Context) {
	k := &s.lk
	s.group(c, "Licenses", func() {
		for i, n := range append(s.svc.System.Notices(), fontNotices()...) {
			open := s.licenseOpen == i
			// A license's text shows with its row.
			s.found(c, n.Name, "", func() {
				s.row(c, n.Name, "", func() {
					label := "show"
					if open {
						label = "hide"
					}
					if s.button(c, "license-"+count(i), label, false).Clicked() {
						if open {
							s.licenseOpen = -1
						} else {
							s.licenseOpen = i
						}
					}
				})
				if open {
					ui.Text(c, n.Text).FontSize(k.captionSize()).Font(monospace()).TextColor(k.muted).Padding(8, 0).Selectable()
				}
			})
		}
	})
}
