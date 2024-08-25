mod consts;
use consts::{NAME, CONFIG_NIX, HELP};
use std::{
    fs, 
    process::{self, Command}, 
    env
};

fn help_msg() -> String {
    let version = env!("CARGO_PKG_VERSION");
    let header = format!("NAPT v{version} (Nix Advanced Package Tool)");
    format!(" {}\n\n{}\n{}", NAME.trim(), header, HELP)
}

fn install(args: &Vec<String>) {
    let text = match fs::read_to_string(CONFIG_NIX) {
        Ok(t) => t,
        Err(e) => {
            println!("Error opening `{}`: {:?}", CONFIG_NIX, e);
            process::exit(0);
        }
    };
    let packages_line = "environment.systemPackages";
    
    let mut new_conf = vec![];
    for line in text.split('\n') {
        new_conf.push(line.to_string());
        if line.contains(packages_line) {
            for package in args.iter().skip(2) {
                new_conf.push(format!("\t\t{}", package));
            }
        }
    }

    let text = new_conf.join("\n");

    if let Err(e) = fs::write(CONFIG_NIX, &text) {
        println!("Error writing to `{}`: {:?}", CONFIG_NIX, e);
        process::exit(0);
    }
}

fn remove(args: &Vec<String>) {
    let r_packages: Vec<&str> = args
        .iter()
        .skip(2)
        .map(|s| s.as_str())
        .collect();

    let text = match fs::read_to_string(CONFIG_NIX) {
        Ok(t) => t,
        Err(e) => {
            println!("Error opening `{}`: {:?}", CONFIG_NIX, e);
            process::exit(0);
        }
    };
    let packages_line = "environment.systemPackages";
    
    let mut new_conf = vec![];
    let mut scan_package = false;

    for line in text.split('\n') {
        if line.contains(packages_line) && !scan_package {
            scan_package = true;
        }
        else if scan_package && line.trim() == "];" {
            scan_package = false;
        }
        else if scan_package && r_packages.contains(&line.trim()) {
            continue;
        }

        new_conf.push(line.to_string());
    }

    let text = new_conf.join("\n");

    if let Err(e) = fs::write(CONFIG_NIX, &text) {
        println!("Error writing to `{}`: {:?}", CONFIG_NIX, e);
        process::exit(0);
    }
}

fn list() {
    let text = match fs::read_to_string(CONFIG_NIX) {
        Ok(t) => t,
        Err(e) => {
            println!("Error opening `{}`: {:?}", CONFIG_NIX, e);
            process::exit(0);
        }
    };
    let packages_line = "environment.systemPackages";
    
    let mut scan_package = false;

    println!("Installed Packages\n=======================");

    for line in text.split('\n') {
        if line.contains(packages_line) && !scan_package {
            scan_package = true;
        }
        else if scan_package && line.trim() == "];" {
            scan_package = false;
        }
        else if scan_package {
            println!("{}", line.trim())
        }

    }


}

fn nixos_rebuild() {
    // sudo nixos-rebuild switch
    let child = Command::new("nixos-rebuild")
        .arg("switch")
        .spawn();

    match child {
        Ok(mut child) => {
            child.wait().unwrap();
        }
        Err(e) => {
            println!("Error rebuilding nix conf:\n{:?}", e);
            process::exit(0);
        }
    }

}

fn main() {

    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("Error: MISSING ARGUMENTS\n\n{}", HELP);
        process::exit(0);
    }

    match args[1].as_str() {
        "install" => {
            install(&args);
            nixos_rebuild();
        },
        "remove" => {
            remove(&args);
            nixos_rebuild();
        },
        "list" => {
            list();
        },
        "--version" | "version" => {
            let version = env!("CARGO_PKG_VERSION");
            println!("napt v{}", version);
        },
        "--help" | "help" => {
            println!("{}", help_msg());
        }
        _ => {
            println!("Error: UNKNOWN COMMAND `{}`\n\n{}", args[1], HELP);
        },
    }

}
