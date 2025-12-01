# See https://developers.google.com/idx/guides/customize-idx-env
{ pkgs, ... }: {
  # Enable the Docker daemon
  services.docker.enable = true;

  # Add the docker-compose package to the environment
  environment.systemPackages = [
    pkgs.docker-compose
  ];
}
