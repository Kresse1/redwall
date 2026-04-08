package wallpaper

import (
	"fmt"
	"os/exec"
	"strings"
)

type KDESetter struct {
}

func NewKDESetter() *KDESetter {
	return &KDESetter{}
}
func (k *KDESetter) Current() (string, error) {
	var jsSnippet = `                
  var allDesktops = desktops();                                                               
  d = allDesktops[0];                                                                         
  d.currentConfigGroup = Array("Wallpaper", "org.kde.image", "General");                      
  print(d.readConfig("Image"));                                                               
  `

	output, err := exec.Command("qdbus6", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", jsSnippet).Output()
	if err != nil {
		return "", err
	}

	trim := strings.TrimSpace(string(output))
	trim = strings.TrimPrefix(trim, "file://")
	return trim, nil

}

func (k *KDESetter) Set(imagePath string) error {
	var jsSnippet = fmt.Sprintf(`                
  var allDesktops = desktops();                                                               
  for (i=0;i<allDesktops.length;i++) {                                                       
      d = allDesktops[i];                                                                     
      d.wallpaperPlugin = "org.kde.image";
      d.currentConfigGroup = Array("Wallpaper", "org.kde.image", "General");                  
      d.writeConfig("Image", "file://%s");
  }                                                                                           
    `, imagePath)

	err := exec.Command("qdbus6", "org.kde.plasmashell", "/PlasmaShell", "org.kde.PlasmaShell.evaluateScript", jsSnippet).Run()
	if err != nil {
		return err
	}
	return nil
}
