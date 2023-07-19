[![Go](https://github.com/amb1s1/goNetConfig/actions/workflows/go.yml/badge.svg)](https://github.com/amb1s1/goNetConfig/actions/workflows/go.yml)
# GoNetConfig

GoNetConfig is a Go-based application that provides network configuration generation services. It uses gRPC for communication and supports multiple configuration feature generators.

## Installation

To install and run the GoNetConfig application, follow these steps:

1. Make sure you have Go installed on your system.

2. Clone the repository:
   ```
   git clone https://github.com/amb1s1/gonetconfig.git
   ```

3. Navigate to the project directory:
   ```
   cd gonetconfig
   ```

4. Install the dependencies:
   ```
   go mod download
   ```

5. Build the application:
   ```
   go build -o gonetconfig .
   ```

6. Run the application:
   ```
   ./gonetconfig
   ```

   The application will start and listen for gRPC requests on `localhost:50051`.

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

3. Implement the `GeneratorInt` interface defined in the `base/base.go` file. This interface requires the implementation of the `Render` method, which generates the network configuration feature. You can use the existing generators as examples.

4. Register your generator in the `service/service.go` file. Add the import statement for your generator package and include a new function in the `allGenerator` slice that returns an instance of your generator.

   ```go
   allGenerator = []func() base.GeneratorInt{
       aaa.New,
       dns.New,
   }
   ```

5. Implement the necessary logic in your generator to generate the desired network configuration feature. You can use templates or any other method suitable for your feature generation.

6. Update the `proto/service.proto` file to include the necessary message types and enums for your new feature. Refer to the existing messages as examples.

7. Run the following command to regenerate the gRPC code based on the updated `.proto` file:

   ```bash
   protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/service.proto
   ```

8. Rebuild and run the GoNetConfig application to apply your changes.

## Project Structure

The GoNetConfig project follows the following structure:

- `main.go`: The entry point of the application that starts the gRPC server and registers the configuration generators.

- `service/service.go`: Contains the gRPC service implementation and the registry of configuration generators.

- `proto/service.proto`: Defines the gRPC service and message types using Protocol Buffers.

- `features/aaa/aaa.go`: Implements the `aaa` configuration generator.

- `base/base.go`: Defines the base structures and interfaces used by configuration generators.

## Contributing

Contributions to the GoNetConfig project are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the project's GitHub repository.

## License

The GoNetConfig project is licensed under the [MIT License](https://opensource.org/licenses/MIT). See the `LICENSE` file for more details.
