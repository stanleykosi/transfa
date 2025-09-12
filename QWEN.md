The Transfa App Technical Specification

## 1. Project Overview

### 1.1 Application Description
Transfa is a native iOS payment application designed to provide a simple, powerful, and unified money movement experience. It abstracts the complexities of the banking network to offer seamless peer-to-peer payments (via username), self-transfers (withdrawals), detailed payment requests, and innovative "Money Drops." The app also provides users with deep financial analytics to understand their spending habits. It serves both individual ("Personal") and merchant ("Business") customers. The application will leverage the Anchor Banking-as-a-Service (BaaS) platform to intelligently route all financial transfers.

### 1.2 Target Platforms
- **Platform:** iOS
- **Minimum Target Version:** iOS 15.0
- **Primary Framework:** SwiftUI
- **Primary Language:** Swift 5.7+

### 1.3 Core Functionality (MVP)
The Minimum Viable Product (MVP) must include the following features:
- **User Onboarding & Account Setup:** Secure sign-up/login via Clerk, KYC/KYB data collection, and automatic creation of a primary in-app wallet via Anchor.
- **Subscription Model:** A freemium model with a paid tier for premium features (unlimited external transfers, multiple beneficiaries).
- **Intelligent Transfer Routing:** Backend logic that automatically routes payments to a user's internal wallet or external bank account based on their subscription status.
- **Core Money Movement:** A central UI for "Pay Someone Else" (P2P) and "Self Transfer" (withdrawal).
- **Wallet Funding:** A simple interface for users to find and copy their virtual account number to fund their wallet.
- **Beneficiary Management:** Ability for users to add, verify, and manage external bank accounts.
- **Payment Requests:** Creation and management of payment requests with optional descriptions and images, shareable via QR code and link.
- **Money Drop:** A wizard-based feature to create a time-bound, multi-claimant fund pool, shareable via a single QR code/link for internal claims.
- **QR Code Payments:** Static, user-specific QR codes for receiving payments.
- **Financial Analytics:** A dedicated tab showing "Money In vs. Money Out" and categorized spending insights.
- **User & Merchant Profiles:** Differentiated profiles with specific permissions (e.g., "Receive Only" for merchants) and filtered transaction history.
- **Security:** Biometric (Face ID/Touch ID) or PIN approval for all outbound money movements.

### 1.4 Success Criteria
The project is considered successful when all MVP features are implemented, stable, and secure according to this specification. The final application must pass all E2E test scenarios defined in Section 11 and meet the performance benchmarks in Section 12.

---

## 2. Technical Architecture

### 2.1 System Design
The system employs a cloud-native microservices architecture deployed on AWS using Docker and Kubernetes (EKS). Communication between the iOS client and the backend is facilitated by a central API Gateway. Asynchronous tasks and inter-service communication are handled by a message broker (RabbitMQ).

### 2.2 Component Relationships
```mermaid
graph TD
    subgraph "Client"
        A[iOS App (SwiftUI)]
    end

    subgraph "Third-Party Services"
        B[Clerk Auth]
        C[Supabase (PostgreSQL + Storage)]
        D[Anchor BaaS]
    end

    subgraph "AWS Cloud Infrastructure"
        E[API Gateway (Kong/Tyk)] --> F{Backend Microservices}
        F --> G[RabbitMQ]
        F --> H[Redis Cache]
        F --> C
        F --> D
        I[Notification Service] --> G
    end

    A -- REST API Calls --> E
    A -- Authenticates With --> B
    E -- Validates JWT --> B
    A -- Uploads Images --> C
```

### 2.3 Data Flow Diagrams

#### 2.3.1 User Onboarding Flow
```mermaid
sequenceDiagram
    participant iOS_App
    participant Clerk
    participant API_Gateway
    participant Auth_Service
    participant RabbitMQ
    participant Customer_Service
    participant Anchor
    participant Notification_Service
    participant Account_Service

    iOS_App->>Clerk: User signs up/in
    Clerk-->>iOS_App: Returns JWT
    iOS_App->>API_Gateway: POST /onboarding (with JWT, username, KYC data)
    API_Gateway->>Clerk: Validate JWT
    API_Gateway->>Auth_Service: Forward Request
    Auth_Service->>Auth_Service: Create user in Supabase 'users' table
    Auth_Service->>RabbitMQ: Publish [user.created] event
    
    Customer_Service->>RabbitMQ: Consume [user.created] event
    Customer_Service->>Anchor: POST /api/v1/customers (Create Customer)
    Anchor-->>Customer_Service: anchor_customer_id
    Customer_Service->>Customer_Service: Update 'users' table with anchor_customer_id
    Customer_Service->>Anchor: POST /customers/{id}/verification (Trigger KYC)
    
    Anchor-->>Notification_Service: Webhook [customer.identification.approved]
    Notification_Service->>RabbitMQ: Publish [customer.verified] event
    
    Account_Service->>RabbitMQ: Consume [customer.verified] event
    Account_Service->>Anchor: POST /api/v1/accounts (Create DepositAccount)
    Anchor-->>Account_Service: anchor_account_id
    Account_Service->>Account_Service: Create record in 'accounts' table
    
    Anchor-->>Notification_Service: Webhook [account.opened]
    Notification_Service->>iOS_App: Push Notification "Wallet Ready"
```

#### 2.3.2 P2P Transfer Routing Flow
```mermaid
sequenceDiagram
    participant iOS_App
    participant API_Gateway
    participant Transaction_Service
    participant Customer_Service
    participant Subscription_Service
    participant Anchor
    
    iOS_App->>API_Gateway: POST /transactions/p2p (to: @recipient, amount)
    API_Gateway->>Transaction_Service: Forward Request
    
    Transaction_Service->>Customer_Service: GET /users/internal?username=@recipient
    Customer_Service-->>Transaction_Service: Recipient User ID & Details
    
    Transaction_Service->>Subscription_Service: GET /subscriptions/status?user_id={id}
    Subscription_Service-->>Transaction_Service: Subscription status & free transfers used
    
    Transaction_Service->>Transaction_Service: is_eligible_for_external = (sub.is_active OR sub.transfers_used < 5)
    
    alt is_eligible_for_external is TRUE
        Transaction_Service->>Customer_Service: GET /beneficiaries/default?user_id={id}
        Customer_Service-->>Transaction_Service: Default CounterParty ID
        Transaction_Service->>Anchor: POST /transfers (type: NIPTransfer)
        Transaction_Service->>Subscription_Service: POST /subscriptions/increment-usage (if not subscribed)
    else is_eligible_for_external is FALSE
        Transaction_Service->>Transaction_Service: Get Recipient's internal DepositAccount ID
        Transaction_Service->>Anchor: POST /transfers (type: BookTransfer)
        Transaction_Service->>RabbitMQ: Publish [transfer.rerouted.internal] event
    end
```

### 2.4 Architectural Patterns
- **iOS (Frontend):** Model-View-ViewModel (MVVM) will be strictly enforced.
- **Backend (Go Microservices):** A hexagonal/clean architecture pattern will be used.
  - `domain`: Core business models and interfaces, no external dependencies.
  - `app`: Business logic implementation, orchestrates domain objects.
  - `store`: Data persistence logic (repositories).
  - `api`: HTTP handlers, routing, and transport layer concerns.

---

## 3. Technology Stack

| Category         | Technology                 | Version/Specification                                     |
| ---------------- | -------------------------- | --------------------------------------------------------- |
| **Frontend**     | iOS Target                 | 15.0+                                                     |
|                  | Swift                      | 5.7+                                                      |
|                  | UI Framework               | SwiftUI                                                   |
| **Authentication** | User Authentication        | Clerk                                                     |
| **Backend**      | Language                   | Go (Golang)                                               |
|                  | Web Framework              | Chi (recommended)                                         |
|                  | Database Driver            | pgx/v5                                                    |
|                  | Env Configuration          | Viper                                                     |
| **Database**     | Primary Database           | Supabase (PostgreSQL 15+)                                 |
|                  | Storage                    | Supabase Storage                                          |
| **Infrastructure** | Cloud Provider             | Amazon Web Services (AWS)                                 |
|                  | Containerization           | Docker                                                    |
|                  | Orchestration              | Kubernetes (AWS EKS)                                      |
|                  | API Gateway                | Kong or Tyk                                               |
|                  | Caching                    | Redis                                                     |
|                  | Messaging                  | RabbitMQ                                                  |
|                  | Observability              | Prometheus (Monitoring), Grafana (Dashboards), Loki (Logs)|
| **BaaS**         | Banking Provider           | Anchor                                                    |

---

## 4. Development Environment Setup

### 4.1 Prerequisites
- **Hardware:** Apple Mac with Apple Silicon (M1/M2/M3) recommended.
- **Software:**
  - Xcode 14.0 or newer
  - Homebrew
  - Docker Desktop
  - Go 1.21+ (`brew install go`)
  - Supabase CLI (`brew install supabase/tap/supabase-cli`)
  - Git

### 4.2 Third-Party Service Configuration
1.  **Clerk:**
    - Create an account at clerk.dev.
    - Create a new application.
    - Note the **Publishable Key** and **Secret Key**.
    - Navigate to JWT Templates and create a new template. Ensure it includes user metadata.
2.  **Supabase:**
    - Create a new project at database.new.
    - Note the **Project URL**, **anon (public) key**, and **service_role key**.
    - In the SQL Editor, create a new query to run the migration scripts from Section 7.
3.  **Anchor:**
    - Request access to the Anchor Sandbox environment.
    - Once approved, log in to the Sandbox Dashboard.
    - Navigate to Developers -> API Keys and create a new API key. Note this key.
    - Navigate to Developers -> Webhooks. You will configure this later once the Notification service is deployed.

### 4.3 Local Backend Setup
1.  Clone the repository.
2.  In the root of the repository, create a `.env` file from the provided `.env.example`.
3.  Populate `.env` with the keys from the services configured above.
    ```env
    # Supabase
    DATABASE_URL="postgres://postgres:[YOUR-PASSWORD]@[YOUR-DB-HOST]:5432/postgres"
    SUPABASE_SERVICE_KEY="your_supabase_service_role_key"

    # Clerk
    CLERK_SECRET_KEY="your_clerk_secret_key"
    CLERK_JWKS_URL="https://your-clerk-instance/.well-known/jwks.json"

    # Anchor
    ANCHOR_API_KEY="your_anchor_sandbox_api_key"
    ANCHOR_BASE_URL="https://api.sandbox.getanchor.co"

    # RabbitMQ
    RABBITMQ_URL="amqp://guest:guest@rabbitmq:5672/"
    ```
4.  Run `docker-compose up -d --build` to start all microservices and dependencies (RabbitMQ, Redis).

### 4.4 iOS Frontend Setup
1.  Navigate to the `ios/` directory.
2.  Install SwiftLint: `brew install swiftlint`.
3.  Open `Transfa.xcodeproj` in Xcode.
4.  Create a `Configuration.xcconfig` file to store environment variables like the API Gateway URL and Clerk Publishable Key. Do not commit keys to Git.
5.  Select a simulator (e.g., iPhone 14 Pro) and run the application (Cmd+R).

---

## 5. Code Structure and Organization

### 5.1 iOS Application (MVVM)
```
Transfa/
├── App/
│   └── TransfaApp.swift         // App entry point, environment setup
├── Core/
│   ├── Models/                  // Codable data models (e.g., UserProfile.swift)
│   ├── Services/                // Networking layer (e.g., APIService.swift, TransactionService.swift)
│   └── Utilities/               // Formatters, extensions, etc.
├── Features/                    // Each distinct feature or screen
│   ├── Onboarding/
│   │   ├── Views/               // SwiftUI Views (e.g., WelcomeView.swift)
│   │   └── ViewModels/          // ObservableObject classes (e.g., WelcomeViewModel.swift)
│   ├── Home/
│   └── ... (other features)
└── Shared/
    ├── Components/              // Reusable SwiftUI views (e.g., PrimaryButton.swift)
    └── Assets.xcassets/         // Colors, icons, images
```
- **File Naming:** `PascalCase` for all Swift files, matching the primary type within.

### 5.2 Backend Microservice (Go)
```
service-name/
├── cmd/
│   └── main.go                  // Service entry point, config loading, server startup
├── internal/
│   ├── api/
│   │   ├── handlers.go          // HTTP request handlers and response writers
│   │   └── router.go            // Route definitions (e.g., using Chi router)
│   ├── app/
│   │   ├── service.go           // Core business logic, orchestrating domain and store
│   │   └── interfaces.go        // Interfaces for dependencies (e.g., Repository, AnchorClient)
│   ├── domain/
│   │   └── models.go            // Core data structures with JSON and DB tags
│   └── store/
│       └── postgres/
│           └── repository.go    // Database interactions (SQL queries via pgx)
├── pkg/                         // Shared libraries (e.g., custom RabbitMQ publisher)
├── api/
│   └── openapi.yml              // OpenAPI 3.0 specification for this service
├── Dockerfile
└── README.md
```
- **Naming Conventions:** `PascalCase` for Go types, `camelCase` for functions and variables. JSON keys and database columns must be `snake_case`.

---

## 6. API Specifications

All communication from the iOS client goes through the API Gateway, which validates the Clerk JWT and routes to the appropriate service.

### 6.1 Auth Service
- **Endpoint:** `POST /onboarding`
  - **Description:** Completes user registration after Clerk sign-up. Creates the user in the local DB and triggers the async onboarding flow.
  - **Auth:** Clerk JWT required.
  - **Request Body:**
    ```json
    {
      "username": "string (unique, alphanumeric, 3-20 chars)",
      "account_type": "'personal' or 'merchant'",
      "kyc_details": { // For 'personal'
        "full_name": "string",
        "bvn": "string (11 digits)",
        "date_of_birth": "string (YYYY-MM-DD)"
      },
      "kyb_details": { // For 'merchant'
        "business_name": "string",
        "rc_number": "string"
      }
    }
    ```
  - **Success Response (201 Created):**
    ```json
    {
      "user_id": "uuid",
      "username": "string",
      "status": "onboarding_initiated"
    }
    ```
  - **Error Responses:**
    - `400 Bad Request`: Validation failed.
    - `409 Conflict`: Username already exists.

### 6.2 Customer Service
- **Endpoint:** `POST /beneficiaries`
  - **Description:** Adds a new external bank account for the user.
  - **Auth:** Clerk JWT required.
  - **Request Body:**
    ```json
    {
      "account_number": "string",
      "bank_code": "string"
    }
    ```
  - **Success Response (201 Created):** The created beneficiary object.
- **Endpoint:** `GET /beneficiaries`
- **Endpoint:** `DELETE /beneficiaries/{id}`
- **Endpoint:** `GET /users/{username}`
  - **Description:** Fetches a user's public profile and the transaction history *between the viewer and the profile owner*.
  - **Auth:** Clerk JWT required.
  - **Success Response (200 OK):**
    ```json
    {
      "username": "string",
      "profile_image_url": "string",
      "account_type": "string",
      "transaction_history": [ /* Array of transaction objects */ ]
    }
    ```
- **Endpoint:** `PUT /users/me/settings`

### 6.3 Transaction Service
- **Endpoint:** `POST /transactions/p2p`
  - **Description:** Initiates a peer-to-peer transfer.
  - **Auth:** Clerk JWT required.
  - **Request Body:**
    ```json
    {
      "recipient_username": "string",
      "amount": "integer (in kobo)",
      "description": "string (optional)",
      "payment_request_id": "uuid (optional)"
    }
    ```
- **Endpoint:** `POST /transactions/self-transfer`
  - **Description:** Initiates a withdrawal to a linked beneficiary.
  - **Auth:** Clerk JWT required.
  - **Request Body:**
    ```json
    {
      "beneficiary_id": "uuid",
      "amount": "integer (in kobo)"
    }
    ```
- **Endpoint:** `POST /money-drops`
- **Endpoint:** `POST /money-drops/{id}/claim`
- **Endpoint:** `POST /payment-requests`

### 6.4 Asynchronous Events (RabbitMQ)
- **`user.created`**: Published by Auth Service. Consumed by Customer Service.
  ```json
  { "user_id": "uuid", "account_type": "...", "kyc_details": { ... } }
  ```
- **`customer.verified`**: Published by Notification Service. Consumed by Account Service.
  ```json
  { "user_id": "uuid", "anchor_customer_id": "string" }
  ```
- **`transaction.completed`**: Published by Transaction Service. Consumed by Analytics Service.
  ```json
  { "transaction_id": "uuid", "sender_id": "uuid", "recipient_id": "uuid", "amount": "integer", "type": "...", "timestamp": "iso8601" }
  ```

---

## 7. Database Schema

All primary keys are `UUID` and default to `gen_random_uuid()`. All tables include `created_at` and `updated_at` timestamps.

### 7.1 Migration Script: `0000_initial_schema.sql`
```sql
-- Create users table
CREATE TABLE public.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clerk_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    account_type TEXT NOT NULL CHECK (account_type IN ('personal', 'merchant')),
    anchor_customer_id TEXT UNIQUE,
    kyc_status TEXT NOT NULL DEFAULT 'pending' CHECK (kyc_status IN ('pending', 'approved', 'rejected')),
    profile_image_url TEXT,
    allow_sending BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT username_length CHECK (char_length(username) >= 3 AND char_length(username) <= 20)
);

-- Link to auth.users
ALTER TABLE public.users ADD CONSTRAINT fk_auth_users FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE;

-- Other tables (accounts, beneficiaries, transactions, etc.) as defined in the technical spec...
-- (Full SQL for all tables would be included here)
```

### 7.2 RLS Policies Script: `0001_rls_policies.sql`
```sql
-- Enable RLS for all tables
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
-- ... for all other tables

-- Policies for 'users' table
CREATE POLICY "Users can view their own profile." ON public.users FOR SELECT USING (auth.uid() = id);
CREATE POLICY "Users can update their own profile." ON public.users FOR UPDATE USING (auth.uid() = id);

-- Policies for 'transactions' table
CREATE POLICY "Users can view their own transactions." ON public.transactions FOR SELECT
USING (auth.uid() = sender_user_id OR auth.uid() = recipient_user_id);

-- (RLS policies for all other tables would be included here)
```

### 7.3 Storage Setup Script: `0002_storage_setup.sql`
```sql
-- Create storage bucket for user content
INSERT INTO storage.buckets (id, name, public) VALUES ('user_content', 'user_content', true)
ON CONFLICT (id) DO NOTHING;

-- Policies for Storage
CREATE POLICY "Authenticated users can upload to their own folder"
ON storage.objects FOR INSERT TO authenticated
WITH CHECK (bucket_id = 'user_content' AND (storage.foldername(name))[1] = auth.uid()::text);

CREATE POLICY "Content is publicly readable"
ON storage.objects FOR SELECT
USING (bucket_id = 'user_content');
```

---

## 8. UI/UX Implementation

### 8.1 Design System
- **Color Palette:** All colors must be defined in `Assets.xcassets` to support light/dark mode.
  - `Primary`: #4F46E5 (Indigo)
  - `Secondary`: #10B981 (Emerald)
  - `Error`: #EF4444 (Red-500)
  - `Background`: System Background Color
  - `TextPrimary`: Label Color
- **Typography:** Use standard SwiftUI dynamic type styles.
  - `Display`: `.largeTitle` (Bold)
  - `Title`: `.title` (SemiBold)
  - `Body`: `.body` (Regular)
- **Spacing:** Use a 4pt grid system. All paddings and spacings must be multiples of 4.

### 8.2 Core SwiftUI Components
- **`PrimaryButton.swift`**:
  ```swift
  struct PrimaryButton: View {
      let title: String
      let action: () -> Void

      var body: some View {
          Button(action: action) {
              Text(title)
                  .font(.headline)
                  .foregroundColor(.white)
                  .frame(maxWidth: .infinity)
                  .padding()
                  .background(Color("Primary"))
                  .cornerRadius(12)
          }
      }
  }
  ```
- **`FormField.swift`**:
  ```swift
  struct FormField: View {
      let title: String
      @Binding var text: String
      var placeholder: String = ""
      var keyboardType: UIKeyboardType = .default

      var body: some View {
          VStack(alignment: .leading) {
              Text(title).font(.caption)
              TextField(placeholder, text: $text)
                  .keyboardType(keyboardType)
                  .padding()
                  .background(Color(UIColor.systemGray6))
                  .cornerRadius(8)
          }
      }
  }
  ```

---

## 9. Business Logic Rules

### 9.1 Transfer Routing
- **Location:** `transaction-service/internal/app/service.go`
- **Logic:**
  1. Receive P2P transfer request for recipient `R`.
  2. Fetch `R`'s user data and subscription status (`S`).
  3. Evaluate eligibility: `is_eligible = S.status == 'active' || S.monthly_external_transfers_used < 5`.
  4. **IF `is_eligible` is TRUE**:
     a. Fetch `R`'s default `beneficiary` (Anchor CounterParty ID).
     b. Initiate Anchor `NIPTransfer`.
     c. If `S.status == 'free'`, atomically increment `monthly_external_transfers_used`.
  5. **IF `is_eligible` is FALSE**:
     a. Fetch `R`'s main `account` (Anchor DepositAccount ID).
     b. Initiate Anchor `BookTransfer`.
     c. Publish `transfer.rerouted.internal` event to RabbitMQ.

### 9.2 Money Drop Claim Validation
- **Location:** `transaction-service/internal/store/repository.go`
- **Logic:** Must be a single atomic database transaction.
  ```sql
  -- Pseudocode for claim validation
  BEGIN;
  SELECT * FROM money_drops WHERE id = {drop_id} FOR UPDATE;
  -- Check if drop is active, not expired, and claims_made < total_claims
  -- Check if claimant_user_id is NOT in money_drop_claims for this drop_id
  IF all_checks_pass THEN
      INSERT INTO money_drop_claims (...);
      UPDATE money_drops SET claims_made_count = claims_made_count + 1 WHERE id = {drop_id};
  ELSE
      ROLLBACK; -- Abort with error
  END IF;
  COMMIT;
  ```

---

## 10. Security Requirements

- **Authentication:** All API endpoints (except Clerk webhooks) must be protected. The API Gateway is responsible for validating the Clerk-issued JWT on every incoming request.
- **Authorization:** Row Level Security (RLS) on Supabase is the primary authorization mechanism and must be enabled on all tables with user-specific data. Additional checks (e.g., merchant `allow_sending` flag) must be implemented in the service layer.
- **Data Handling:** No Personally Identifiable Information (PII) or sensitive financial data is ever to be written to logs. This includes names, emails, phone numbers, BVNs, and account numbers.
- **Secrets Management:** All secrets (DB connection strings, API keys) must be loaded from environment variables. In production, these will be injected from AWS Secrets Manager.
- **Client-Side Security:**
  - All sensitive data on the iOS device (JWTs, user PIN) **must** be stored in the iOS Keychain. `UserDefaults` is forbidden for this purpose.
  - **Biometric Approval:** Any user action that initiates an outbound transfer (`P2P`, `Self Transfer`, `Money Drop Creation`) must first be approved via a `LocalAuthentication` prompt. The API call is only made after successful biometric validation.

---

## 11. Testing Strategy

### 11.1 Unit Tests
- A minimum of 70% code coverage is required for business logic in Go `app` packages and iOS `ViewModels`.
- **Backend Example:** `TestP2PTransferRouting_FreeUser_ReroutesAfterFiveTransfers`.
- **iOS Example:** `testHomeViewModel_FetchData_SetsIsLoadingCorrectly`.

### 11.2 End-to-End (E2E) Test Scenarios (`XCUITest`)
The following critical user flows must be automated:
1.  **Full Onboarding:** Sign up, complete KYC, and verify the home screen shows a zero balance.
2.  **P2P Transfer (External):** A subscribed user receives a payment and it is routed externally.
3.  **P2P Transfer (Internal Reroute):** A free user with no free transfers left receives a payment, and it is deposited into their in-app wallet.
4.  **Money Drop Full Cycle:** User A creates a drop. User B and C claim it. Wait for expiry and verify remaining funds are returned to User A.
5.  **Payment Request Flow:** User A creates a request with an image. User B pays via the link. Verify request status becomes "Fulfilled".
6.  **Merchant Send Block:** A merchant in "Receive Only" mode attempts a P2P transfer and is blocked by the UI and API.
7.  **Biometric Approval Flow:** Initiate a transfer, verify the biometric prompt appears, cancel it, and confirm the transfer is not sent. Repeat, approve, and confirm the transfer is sent.

---

## 12. Performance Standards

- **API Latency:** 99% of all API calls must complete in under 500ms. P95 should be under 200ms.
- **Database Queries:** No single query should take longer than 50ms. Use indexes on all foreign keys and frequently queried columns.
- **App Performance:** UI animations should maintain 60 FPS. Cold app start time must be under 2 seconds.
- **Asynchronicity:** Any operation that involves external API calls (e.g., to Anchor) or heavy computation must not block the main request-response cycle. Use RabbitMQ for such tasks.

---

## 13. Error Handling

### 13.1 Backend Error Responses
All 4xx/5xx errors must return a standardized JSON:API error object.
```json
{
  "errors": [
    {
      "status": "400",
      "title": "Bad Request",
      "detail": "Username must be alphanumeric."
    }
  ]
}
```
Go services must handle all errors gracefully and never `panic`.

### 13.2 Frontend Error Handling
- The central `APIService` must handle all network and decoding errors.
- ViewModels must expose an `@Published` property (e.g., `errorMessage: String?`).
- Views must observe this property and display a user-friendly alert or inline message. The app must never crash due to an unhandled API error.

---

## 14. Deployment Configuration

### 14.1 CI/CD
- A CI/CD pipeline (e.g., GitHub Actions) must be configured.
- On every push to a `feature/*` branch, it must run linters and unit tests.
- On merge to `develop`, it must build Docker images, push them to ECR, and deploy to a staging EKS cluster.
- On merge from `develop` to `main`, a manual approval step is required to deploy to the production EKS cluster.

### 14.2 Environment Variables
Each microservice will require the following environment variables (at minimum):
- `PORT`
- `DATABASE_URL`
- `RABBITMQ_URL`
- `REDIS_URL`
- `CLERK_JWKS_URL`
- `ANCHOR_API_KEY`
- `ANCHOR_BASE_URL`

### 14.3 iOS Deployment
- The app will be deployed to the Apple App Store.
- The build process will use different schemes for Staging and Production environments, which will point to the respective backend URLs via the `Configuration.xcconfig` file.

---

## 15. Code Quality Standards

- **Formatting:** Code must be automatically formatted by `gofmt` (Go) and `SwiftLint` (iOS) in a pre-commit hook and in the CI pipeline.
- **Documentation:** All public types, functions, and methods in Go must have GoDoc comments. All public methods in iOS Services and ViewModels must have documentation comments explaining their purpose, parameters, and return values.
- **Pull Requests:** All code must be merged into `develop` via a Pull Request. PRs must be small, focused on a single task, and have at least one approval. Squash merges are to be used to maintain a clean Git history.