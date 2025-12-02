# See https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Enable the Docker daemon
  services.docker.enable = true;

  # Add required packages
  packages = [
    pkgs.docker-compose
    pkgs.awscli
    pkgs.go
  ];
}
