# Event oriented application framework written in go
## Focus on 
1. Seperation of Concern 
    - business logic, application logic and low-level logic are abstracted into different layers
        -  Low-level interaction: `InstructionSet`
        -  Application logic: `Endpoint`
        -  Business logic: `Microservice`
3. Declarative and Observability
    - Built-in service discovery via `etcd config`
        - Service microservice are managed by its communication channel (`topic` in kakfa)
    - Each entity are represeneted and managed by its ID
        - `Event` responsibility chain are managed 
        - built-in `logger middleware` push to `elasticsearch` 
        - routing mechanism are defined in `etcd` database and synced in realtime
4. Scalable and Robust 
    - robust:
        - prototype design pattern
            - no additional module abstraction 
        - type flexibility
            - data are passed as key-value pair 
        - prebuild:
            - pre-builded user database (similiar to .NET identity framwork)
    - scalability:
        - fully asynchrous communication design pattern
        - client-server communication use `gRPC bidirectional stream`
        - client authentication and authorization are performed only once
6. Event-driven 
    - `UserRequest` are abstracted into `Event`
        - `Event` are passed and processed between multiple `Microservice`'s `Endpoint`
        - all `Endpoint` operate on it own key-value segment called `Action`
        - all `UserRequest` are translated into `Event` using predefined rule (created by domain developer)

