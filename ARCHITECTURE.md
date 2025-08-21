# Application Architecture

This warehouse management system follows a **clean architecture** pattern with a layered structure to ensure modularity, scalability, and maintainability. The architecture is divided into several key components:

## Layers

1. **Presentation Layer (Controller)**:

   - Handles HTTP requests and responses using the Fiber framework.
   - Contains `ProductController` and `AuthController` to manage API endpoints.
   - Integrates with middleware for authentication (`Authenticate`) and authorization (`Authorize`).

2. **Application Layer (Use Case)**:

   - Contains business logic and orchestrates interactions between repositories and controllers.
   - Implements use cases like `CreateProduct`, `GetProductStocksList`, `TrackStockMovement`, and `GetDashboardSummary`.
   - Ensures validation, error handling, and data transformation (e.g., DTO mapping).

3. **Data Layer (Repository)**:

   - Manages database operations using GORM with PostgreSQL.
   - Provides CRUD operations and complex queries (e.g., joins for dashboard details).
   - Abstracts database interactions for modularity.

4. **Middleware**:
   - `AuthMiddleware`: Handles JWT authentication and sets `userID`/`role` in context.
   - `ProductMiddleware`: Enforces RBAC based on user roles (user, admin, super_admin).
   - Ensures security and access control.

## Data Flow

- **Request**: Enters via Fiber router → Middleware (auth/authorization) → Controller.
- **Processing**: Controller calls Use Case → Use Case interacts with Repository → Database.
- **Response**: Repository returns data → Use Case transforms (if needed) → Controller sends JSON response.

## Key Features

- **Modularity**: Each module (products, warehouse, stock) is self-contained with its own routes, controllers, and use cases.
- **Scalability**: Layered design allows easy addition of new features or services.
- **Security**: RBAC and JWT-based authentication protect endpoints.
- **Audit Trail**: Stock movements are logged for tracking changes.

## Dependencies

- **Fiber**: Web framework for routing and middleware.
- **GORM**: ORM for PostgreSQL interaction.
- **Validator**: Input validation.
- **Logrus**: Logging.

## Future Improvements

- Add unit/integration tests.
- Implement database migrations.
- Introduce caching (e.g., Redis) for dashboard data.
