//go:build windows

package giode

import (
	"errors"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// Explorer's classic context menu reads verbs from these keys under
// HKEY_CURRENT_USER\Software\Classes, which needs no administrator
// rights. The keys created are recorded under Software\<prefix> so
// uninstalling finds them.

const contextActionsRecord = "GiodeContextActions"

func installContextActions(prefix, exe string, actions []ContextAction) error {
	uninstallContextActions(prefix)
	var created []string
	add := func(parent, id, label, arg string) error {
		path := `Software\Classes\` + parent + `\shell\` + prefix + "." + id
		k, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		k.SetStringValue("MUIVerb", label)
		k.SetStringValue("Icon", exe)
		k.Close()
		c, _, err := registry.CreateKey(registry.CURRENT_USER, path+`\command`, registry.ALL_ACCESS)
		if err != nil {
			return err
		}
		defer c.Close()
		created = append(created, path)
		return c.SetStringValue("", `"`+exe+`" `+contextActionFlag+id+` "`+arg+`"`)
	}
	var errs []error
	for _, c := range actions {
		if c.onFiles() {
			if len(c.Extensions) == 0 {
				errs = append(errs, add("*", c.ID, c.Label, "%1"))
			}
			for _, ext := range c.Extensions {
				if !strings.HasPrefix(ext, ".") {
					ext = "." + ext
				}
				errs = append(errs, add(`SystemFileAssociations\`+ext, c.ID, c.Label, "%1"))
			}
		}
		if c.Folders {
			errs = append(errs, add("Directory", c.ID, c.Label, "%1"))
		}
		if c.FolderBackground {
			errs = append(errs, add(`Directory\Background`, c.ID, c.Label, "%V"))
		}
	}
	if k, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\`+prefix, registry.ALL_ACCESS); err == nil {
		k.SetStringsValue(contextActionsRecord, created)
		k.Close()
	}
	return errors.Join(errs...)
}

func uninstallContextActions(prefix string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\`+prefix, registry.ALL_ACCESS)
	if err != nil {
		return nil
	}
	defer k.Close()
	paths, _, _ := k.GetStringsValue(contextActionsRecord)
	var errs []error
	for _, p := range paths {
		if !strings.Contains(p, `\shell\`+prefix+".") {
			continue // not ours
		}
		errs = append(errs, deleteKeyTree(registry.CURRENT_USER, p))
	}
	k.DeleteValue(contextActionsRecord)
	return errors.Join(errs...)
}

// deleteKeyTree deletes a key and its subkeys.
func deleteKeyTree(root registry.Key, path string) error {
	k, err := registry.OpenKey(root, path, registry.ALL_ACCESS)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return err
	}
	subs, _ := k.ReadSubKeyNames(-1)
	k.Close()
	for _, s := range subs {
		if err := deleteKeyTree(root, path+`\`+s); err != nil {
			return err
		}
	}
	return registry.DeleteKey(root, path)
}
