[![Go](https://github.com/amb1s1/goNetConfig/actions/workflows/go.yml/badge.svg)](https://github.com/amb1s1/goNetConfig/actions/workflows/go.yml)
# GoNetConfig

GoNetConfig is a Go-based application that provides network configuration generation services. It uses gRPC for communication and supports multiple configuration feature generators with per-entity (company/tenant) customization.

## Installation

1. Make sure you have Go installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/amb1s1/gonetconfig.git
   ```

3. Navigate to the project directory:
   ```
   cd gonetconfig
   ```

4. Build and run:
   ```
   make build
   make run
   ```

   The application will start and listen for gRPC requests on `localhost:50051`.

## Entity Configuration

GoNetConfig supports multi-tenant configuration through the **entity** system. Each entity (e.g. a company) can have its own set of enabled features and feature-specific parameters.

### Setup

Copy the template config to your home directory:

```bash
make entity-init
```

This creates `~/.gonetconfig/entity.yml` from the template in `configs/entity.yml`.

### Config File Structure

```yaml
version: "1.0"

# Feature parameters per entity
params:
  companyA:
    aaa:
      secret: "companyA-secret-123"
    logger:
      management_interface: "loopback0"
      logger_server_ips:
        - "10.0.1.1"
        - "10.0.1.2"
  companyB:
    aaa:
      secret: "companyB-secret-456"

# Which features each entity enables
features:
  companyA:
    - aaa
    - logger
  companyB:
    - aaa

# Named entity definitions
entity:
  companyA:
    features: companyA
    params: companyA
  companyB:
    features: companyB
    params: companyB

selected_entity: companyA
```

### Config Loading Order

1. `~/.gonetconfig/entity.yml` (user config)
2. `configs/entity.yml` (template/fallback)
3. If neither exists, all features run with hardcoded defaults

### Switching Entities

Edit the `selected_entity` field in your config file to switch between entities. The selected entity controls:

- **Which features are generated** (via the `features` list)
- **What parameter values are used** (via the `params` map)

## Usage

The GoNetConfig application exposes a gRPC service called `GoNetConfigService`, which provides the following RPC method:

### GetConfigGen

This method generates network configurations based on the requested configuration generator type and device information.

```protobuf
rpc GetConfigGen(ConfigGenRequest) returns (ConfigGenResponse);
```

The `ConfigGenRequest` message contains the following fields:

- `config_gen_type`: The configuration generator type. It can be `CGT_FULL` or `CGT_PARTIAL`.

- `device`: Information about the device for which the configuration is generated. It includes the `name`, `vendor`, and `model` fields.

The `ConfigGenResponse` message contains the following fields:

- `status`: The status of the configuration generation process.

- `config_feature`: The generated network configuration feature. It includes the `name`, `version`, and `configuration` fields.

## Configuration Generators

The GoNetConfig application supports multiple configuration generators. Each generator is responsible for generating a specific network configuration feature. To add a new feature generator, follow the steps below:

1. Create a new package for your feature generator under the `features` directory. For example, if you are adding a feature generator for DNS configurations, create a package named `dns` under `features`.

2. Inside the package, create a Go file that contains the implementation of your generator. For example, create a file named `dns.go`.

3. Implement the `GeneratorInt` interface defined in the `base/base.go` file. This interface requires the implementation of the `Render`, `Generators`, and `SetParams` methods.

   ```go
   type GeneratorInt interface {
       Render(context.Context, *pb.ConfigGenRequest, *pb.ConfigGenResponse) *pb.ConfigFeature
       Generators() *Generator
       SetParams(map[string]interface{})
   }
   ```

4. Register your generator in the `service/service.go` file. Add the import statement for your generator package and include a new function in the `allGenerator` slice that returns an instance of your generator.

   ```go
   allGenerator = []func() base.GeneratorInt{
       aaa.New,
       logger.New,
       dns.New,
   }
   ```

5. Add your feature's parameters to `configs/entity.yml` under the appropriate entity params and features lists.

6. Update the `proto/service.proto` file to include the necessary message types and enums for your new feature if needed.

7. Run the following command to regenerate the gRPC code based on the updated `.proto` file:

   ```bash
   make proto
   ```

8. Rebuild and run the GoNetConfig application:

   ```bash
   make build && make run
   ```

## Project Structure

```
.
├── main.go                         # Entry point, loads entity config and starts server
├── Makefile                        # Build, test, and development commands
├── base/
│   └── base.go                     # GeneratorInt interface and base types
├── configs/
│   └── entity.yml                  # Template entity configuration
├── entity/
│   ├── entity.go                   # Entity config loading, selection, parsing
│   └── entity_test.go              # Entity tests
├── features/
│   ├── aaa/
│   │   ├── aaa.go                  # AAA feature generator
│   │   ├── cisco.go                # Cisco AAA renderer
│   │   ├── ciscoxr_template.go     # AAA template
│   │   └── aaa_test.go             # AAA tests
│   └── logger/
│       ├── logger.go               # Logger feature generator
│       ├── cisco.go                # Cisco logger renderer
│       ├── ciscoxr_Template.go     # Logger template
│       └── logger_test.go          # Logger tests
├── service/
│   ├── service.go                  # gRPC service, generator registry, entity integration
│   └── service_test.go             # Service tests
├── server/
│   ├── server.go                   # gRPC server setup
│   └── server_test.go              # Server tests
├── client/
│   └── main.go                     # gRPC client example
└── proto/
    ├── service.proto               # Protobuf service definition
    ├── service.pb.go               # Generated protobuf code
    └── service_grpc.pb.go          # Generated gRPC code
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build the application binary |
| `make run` | Run the application |
| `make test` | Run all tests |
| `make test-v` | Run all tests with verbose output |
| `make vet` | Run `go vet` on all packages |
| `make clean` | Remove build artifacts |
| `make proto` | Regenerate gRPC code from `.proto` files |
| `make entity-init` | Copy template entity config to `~/.gonetconfig/` |
| `make help` | Show all available targets |

## Contributing

Contributions to the GoNetConfig project are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the project's GitHub repository.

## License

The GoNetConfig project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the `LICENSE` file for more details.
