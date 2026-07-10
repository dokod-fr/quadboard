# QuadBoard Configuration

QuadBoard offers a flexible configuration system. You can configure it using a YAML file, environment variables, or a combination of both.

## Priority Order

Configuration is loaded using the following priority (from highest to lowest):

  - `Environment Variables`: Override all other settings.
  - `YAML File`: Overrides default constants.
  - `Default Constants`: Built-in sensible defaults.

## Configuration File

By default, QuadBoard will look for a config.yaml file located next to the QuadBoard executable. 

You can specify a custom path for the configuration file by setting the QUADBOARD_CONFIG_PATH environment variable.
Example config.yaml

```yml
server:
  address: "0.0.0.0:8080"
    read_timeout: 5
      write_timeout: 10
logging:
  level: "info"
  format: "text"
providers:
  quadlet:
    paths:
      - /etc/containers/systemd/
      - /opt/quadboard/quadlets
```

## Environment Variables

All environment variables start with the QUADBOARD_ prefix.

### Server Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|	
| QUADBOARD_CONFIG_PATH	        | Path to a custom YAML configuration file.        |	(Checks for config.yaml next to the binary) |
| QUADBOARD_SERVER_ADDRESS	    | The address and port the HTTP server listens on. |	0.0.0.0:8080 |
| QUADBOARD_SERVER_READ_TIMEOUT	| HTTP read timeout in seconds.                    |	5  |
| QUADBOARD_SERVER_WRITE_TIMEOUT|	HTTP write timeout in seconds.                   | 10  |
  
### Logging Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|
| QUADBOARD_LOGGING_LEVEL	      | Log level (e.g., debug, info, warn, error).	| info |
| QUADBOARD_LOGGING_FORMAT	    | Log format (e.g., text, json).	            | text |
  
### Providers Configuration

| Variable                      | Description               | Default                |
|-------------------------------|---------------------------|------------------------|	
| QUADBOARD_QUADLET_PATHS	      | Comma-separated list of directories to scan for Quadlet files.	| /etc/containers/systemd/, ~/.config/containers/systemd/ |
  

> Note: When using QUADBOARD_QUADLET_PATHS, it completely replaces the default paths and the YAML configuration paths.