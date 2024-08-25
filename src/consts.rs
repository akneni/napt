pub const CONFIG_NIX: &str = "/etc/nixos/configuration.nix";

pub const HELP: &str = "Usage: napt [COMMAND] [OPTIONS]

A simple CLI tool to manage system packages in NixOS by modifying the `configuration.nix` file.

Commands:
  install [PACKAGE...]
        Installs one or more packages by adding them to the `environment.systemPackages` section of the NixOS configuration.
        Example: napt install vim git

  remove [PACKAGE...]
        Removes one or more packages from the `environment.systemPackages` section of the NixOS configuration.
        Example: napt remove vim git

  list
        Lists all the packages currently installed in the `environment.systemPackages` section of the NixOS configuration.
        
  --version
        Displays the version of this tool.

  --help
        Displays this help message.


Options:
  --help       Show this message and exit.
  --version    Displays the version of this tool

Examples:
  1. To install a package (e.g., vim):
     $ napt install vim

  2. To remove a package (e.g., vim):
     $ napt remove vim

  3. To list all installed packages:
     $ napt list

Notes:
  - The `install` and `remove` commands will automatically trigger a NixOS rebuild after modifying the configuration.
  - Ensure you have the necessary permissions to modify `/etc/nixos/configuration.nix`.
";

pub const NAME: &str = r#"
 _   _          _____ _______ 
| \ | |   /\   |  __ \__   __|
|  \| |  /  \  | |__) | | |   
| . ` | / /\ \ |  ___/  | |   
| |\  |/ ____ \| |      | |   
|_| \_/_/    \_\_|      |_|   
                            
                               
"#;