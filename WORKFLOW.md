## SLV Workflow

### User Flow

**Initializing SLV in dev machine with remote profile repository**
 - User initializes SLV by providing a remote profile repository
 - SLV pulls the profile from remote repository and stores it locally under slv app data directory
 - SLV is set to sync from remote profile on a periodic basis

**Creating user environment for the current user**
 - User invokes SLV to register current dev machine as a user environment by providing information about the current user
 - SLV generates a new environment key pair
 - Stores the environment secret key in the OS credential store (such as keychain)
 - Serialises the environment metadata (along with the public key) as Serialised Environment Definition string and returns to the user
 - User shares the serialised environment definition string with the admin and requests to add it to the remote profile

**Creating/Initializing a Project**
 - User invokes SLV to initialize a specific directory as a project directory
 - SLV creates a .slv directory in the specified directory, thereby marking it as a SLV project

**Creating a Vault**
 - User invokes SLV to create a vault by specifying environments to share it with along with the name of the vault (if it is project bound) or path to the vault file.
 - SLV gets the specified environment public keys and root public key from profile
 - SLV generates a new vault key pair and encrypts the vault private key with the specified environment public keys and root public key
 - SLV then creates a vault file with extension `<vault_name>.vault.slv` and writes the vault public key and the vault private key wrappings.

**Adding Secrets to Vault**
 - User invokes SLV to add a secret (as key:value) to the vault. (Need to specify vault name or vault file path)
 - SLV reads the vault public key from the vault file `<vault_name>.vault.slv`
 - SLV encrypts the value of the secret using the vault public key and writes the encrypted secret (key:encrypted(value)) to the same vault file

**Sharing an existing Vault with another environment**
 - User requests the admin to share an existing vault to an environment by specifying the environment information

**Sequence Diagram**
```mermaid
sequenceDiagram
    actor User
    participant KeyChain
    participant SLV
    participant RemoteProfile
    participant UserLocalProfile
    participant LocalDir
    participant Admin

    Note over User,Admin: Initializing SLV in dev machine with remote profile repo (sync from remote)
    User->>SLV: Initialize with profile_repository
    SLV->>RemoteProfile: pull profile from remote git repository
    RemoteProfile-->>UserLocalProfile: 
    SLV-->UserLocalProfile: Sync remote profile with local profile periodically
    Note over User,Admin: Creating user environment for self
    User->>SLV: Create local user environment
    SLV-->>SLV: Generates env keypair & env definition
    SLV->>KeyChain: Env Secret Key
    SLV->>User: Env Defintion (with public key)
    User->>Admin: Requests to update env definition in remote profile
    Note over User,Admin: Creating service environment
    User->>SLV: Create service environment (providing service metadata)
    UserLocalProfile-->>SLV: Gets root public key
    SLV-->>SLV: Generates env keypair & env definition
    SLV-->>SLV: Wraps env secret key with root public key
    SLV->>User: Env Wrapped Key (encrypted by root)
    SLV->>User: Env Defintion (with public key)
    User->>Admin: Requests to update env definition in remote profile and add the env secret key to target environment
    Note over User,Admin: Creating a project
    User->>SLV: Create project in a given directory
    SLV->>LocalDir: Creates .slv directory in the given directory
    Note over User,Admin: Creating a vault
    User->>SLV: Create vault abc with access to env1, env2
    UserLocalProfile-->>SLV: Reads env1, env2 & root public keys
    SLV-->>SLV: Generate vault key pair
    SLV-->>SLV: Encrypts vault secret key with env1, env2 & root public keys
    SLV->>LocalDir: Writes it all to abc.vault.slv file under .slv/vaults
    Note over User,Admin: Adding secrets to vault
    User->>SLV: Add secret to abc_vault - foo:bar
    LocalDir-->>SLV: Reads vault public key from .slv/vaults/abc.vault.slv
    SLV-->>SLV: Encrypts "bar" with vault public key
    SLV->>LocalDir: Writes back foo:encrypted(bar) to .slv/vaults/abc.vault.slv
    Note over User,Admin: Sharing an existing vault to new Environment
    User-->>Admin: Requests admin to share vault with an environment from the remote profile
```

### Admin Flow

**Creating a new profile and setting it to sync to remote (forward sync)**
 - Admin invokes slv to create a new profile
 - Admin inputs the profile name (example_profile) and remote repo (for forward sync)
 - SLV creates a new profile directory (example_profile) and sets the given repository as a remote sync repo

**Creating a root environment in the profile**
 - Admin invokes SLV to create a root environment for the profile
 - SLV creates a root key pair and writes the root public key to the profile into environments.slv file under the profile directory (example_profile)
 - SLV also returns the root secret key back to the admin
 - Admin stores the root secret key in a credential store (such as KMS/secret manager) accessible only by the admin

**Modifying settings in a profile**
 - Admin invokes slv to update settings such as sync_interval, allow_vault_sharing
 - SLV writes the settings to the profile into settings.slv file under the profile directory (example_profile)

**Creating service environments in profile**
 - Admin invokes SLV to create a service environment by specifying the metadata for the service such as name, email, tags
 - SLV generates a new environment keypair for the service environment
 - SLV writes the environment public key and the metadata to the profile into environments.slv file
 - SLV returns the created environment's secret key back to admin
 - Admin writes the environment secret key into the credential store that will be accessible only by the target service

**Adding user created environments to profile**
 - Admin receives request from user to add an environment to remote profile. The user sends the serialised environment definition string and the wrapped environemnt secret key (encrypted by root) along with the request.
 - Admin reviews the request and invokes slv to add the environment to profile.
 - SLV writes the environment public key and the metadata to the profile into environments.slv file
 - Admin invokes SLV to decrypt the wrapped environment secret key
 - SLV reads the root secret key and uses it to decrypt the wrapped environment secret key and returns the decrypted secret key to admin
 - Admin writes the decrypted environment secret key into the credential store that will be accessible only by the target service

**Sync local profile to remote**
 - Admin invokes SLV to sync the profile to a remote repository
 - SLV reads the profile and performs a git commit & push action by updating the changes in the local profile to remote thereby keeping it in sync

**Sharing existing vaults to environments upon request**
 - User requests admin to share a vault (vault_abc) in a given repo to an environment (env_2)
 - Admin validates the request for necessary approval and invokes SLV to share vault_abc with env_2 using root secret key
 - SLV receives root secret key from the credential store that is accessible only by the admin
 - SLV also reads the env_2 metadata and gets the env_2 public key
 - SLV then adds access to env_2 for vault_abc by decrypting the vault secret key with root secret key and re-encrypting it with env_2 public key and writing it as an additional entry in the vault file
 - One it is done, admin sends a PR for the same in the specified repo that has the vault file

**Sequence Diagram**
```mermaid
sequenceDiagram
    actor Admin
    participant SLV
    participant RemoteProfile
    participant AdminLocalProfile
    participant KMS
    participant ProjectRepository
    actor User

    Note over User,Admin: Creating a new profile and setting forward sync (to remote)
    Admin->>SLV: Create a new profile xyz
    SLV->>AdminLocalProfile: Creates a profile under slv app data directory
    Note over User,Admin: Creating a root environment for the profile
    Admin->>SLV: Create root environment
    SLV-->>SLV: Generates a root environment keypair
    SLV->>AdminLocalProfile: Adds the root env public key to profile (environments.slv)
    SLV->>Admin: Returns root env secret key
    Admin->>KMS: Stores the root secret key in a credential store (accessible only by admin)
    Note over User,Admin: Modifying settings in local profile
    Admin->>SLV: Change settings (sync_interval, allow_vault_sharing)
    SLV->>AdminLocalProfile: Updates settings in profile (settings.slv)
    Note over User,Admin: Creating service environments in profile
    Admin->>SLV: Create new service environment
    SLV-->>SLV: Generates a service environment keypair
    SLV->>AdminLocalProfile: Adds env public key and env metadata to profile (environments.slv)
    SLV->>Admin: Returns service env secret key
    Admin->>KMS: Stores the env secret key in a credential store (accessible only by target environment)
    Note over User,Admin: Adding user created environments to profile
    User->>Admin: Requests to add env (providing env defintion string) to remote profile with env secret key wrapped by root env public key
    Admin->>SLV: Add env to profile (env_definition_string)
    SLV->>AdminLocalProfile: Adds env public key and env metadata to profile (environments.slv)
    Admin->>SLV: Get env secret key from wrapped env secret key
    KMS-->>SLV: Gets root secret key from credential store
    SLV-->>SLV: Decrypts wrapped env secret key with root secret key
    SLV-->Admin: Returns decrypted env secret key
    Admin->>KMS: Stores the env secret key in a credential store (accessible only by target environment)
    Note over User,Admin: Sync local profile to remote repo
    Admin->>SLV: Sync to remote
    AdminLocalProfile-->>SLV: Reads local profile
    SLV->>RemoteProfile: Pushes local profile to remote using git
    Note over User,Admin: Processing vault sharing requests
    ProjectRepository-->>User: 
    User->>Admin: Requests to share vault vault_abc with environment env_2
    Admin-->>Admin: Validates the request
    ProjectRepository-->>Admin: 
    Admin->>SLV: Share vault_abc with env_2 using root secret key
    RemoteProfile-->>SLV: Receives env_2 public key
    KMS-->>SLV: Receives root secret key from KMS
    SLV->>ProjectRepository: Adds access to env_2 for vault_abc
    Admin->>ProjectRepository: Sends a PR for the same
```

### Reading Secrets

**Reading secrets from vault**
 - Service reads the secret from vault by specifying the secret name (foo) and the vault name (vault_abc)
 - SLV reads the environment secret key from cloud credential store (for service environemnts) or system keychain (for user machine) or simply from environment variable (SLV_SECRET_KEY)
 - SLV reads the encrypted secret value (value of foo from the vaults file)
 - SLV attempts to unlock the vault with the secret key
 - If the given secret key has access to the vault, SLV decrypts and returns the secret
 - If the given secret key doesn't have access to the vault, SLV returns an error

**Sequence Diagram**
```mermaid
sequenceDiagram
    participant Env as User/Service
    participant SLV
    participant Vault
    participant SecretKey as SecretKey (from KMS/Env Var)

    Note over Env,SecretKey: Reading secret from vault
    Env->>SLV: Read secret with a given name (foo)<br/> from the vault (vault_abc)
    SecretKey-->>SLV: SLV reads the environment secret key <br/> from cloud credential store (for service environemnts) <br/> or system keychain (for user machine) <br/> or simply from environment variable (SLV_SECRET_KEY)
    Vault-->>SLV: SLV reads the encrypted secret value <br/> (value of foo from the vaults file)
    SLV-->>SLV: SLV attempts to unlock the vault with the secret key
    Note over Env,SLV: If the given secret key has access to the vault
    SLV->>Env: Decrypts and returns the secret
    Note over Env,SLV: If the given secret key doesn't have access to the vault
    SLV->>Env: Returns an error that <br/> the secret key does not have access
```