sudo apt autoclean 
sudo apt clean 
sudo apt autoremove
sudo journalctl --vacuum-time=7d
sudo journalctl --vacuum-size=200M
sudo snap list --all | awk '/disabled/{print $1, $3}' | while read pkg revision; do sudo snap remove "$pkg" --revision="$revision"; done
sudo rm -rf /tmp/*
sudo rm -rf /var/tmp/*
sudo apt autoclean 
sudo apt clean 
rm -rf ~/.cache/*