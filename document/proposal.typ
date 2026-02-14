#set document(
  title: [Project proposal: Koopify]
)

#set text(
  font: "IBM Plex Serif",
  size: 12pt
)

#set par(justify: true)

#show link: it => text(
  fill: rgb(40, 80, 120),
)[#it]

#title()

= General description

*Koopify* is an open-source e-commerce platform designed for selling both digital and physical goods. The platform provides core e-commerce functionalities, including product browsing, basic inventory management, user authentication and checkout processing through third-party payment providers.

The demo instance of koopify (#link("https://koopify.piguy.nl")[koopify.piguy.nl]) will simulate an online store, dedicated to trading cards and related merchendise. The system will demonstrate full e-commerce workflows which are possible with the Koopify software from the perspective of both an average user and an administrator.

= List of available functionalities

1. Product catalog:
  - Browse all available products
  - View individual product pages with detailed information
  - Browse curated collections (e.g., "Holiday Specials")
  - Filter, sort and search functionality
  - A local index of previously visited product pages

2. User management:
  - User registration and authentication
  - Account management: Reset password or delete account
  - Role-based access control

3. Administrative features
  - Promote registered users to administrator status
  - Trigger a password reset flow on behalf of any user
    - The selected user recieves an email with password reset instructions
  - Manage products (CRUD)
  - Manage curated collections

4. Basic inventory management
  - Add or remove products from the store
  - Set the available stock for any product

5. Checkout flow
  - Shopping cart functionality
  - Integration with third-party payment providers
  - Order confirmation and transaction recording

#pagebreak()

= Development

The backend follows a feature-based architecture, inspired by the MVC model. Each domain (e.g. user, product, order) encapsulates its own HTTP handlers, DTOs, services and repositories. The Echo web framework will be used for handling requests.

The database in use is PostgreSQL, and pgx is used to communicate with the database. SQLC is used to make development easier, along with golang-migrate to handle database migrations.

The frontend is a simple Vue.js single page application, where routing and UI rendering is driven by javascript instead of plain HTML.

The whole project is containerised. GitHub actions workflows will be used to automatically build a new image for every commit to main, and these images will automatically be deployed to the backend. This provides a full CI/CD workflow for the project.

= Timeline

The development of the project is expected to take around three weeks, with an additional week for testing and hotfixes.

= Links
- Git repository: https://github.com/piguycs/koopify/
- Demo instance: https://koopify.piguy.nl/
- Echo: https://echo.labstack.com/
- SQLC: https://sqlc.dev/
- PostgreSQL: https://www.postgresql.org/
- Vue: https://vuejs.org/
