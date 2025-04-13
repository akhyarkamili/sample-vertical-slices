# Sample Vertical Slices for a Loan App

The system supports proposing new loans, approving new loans, and partially implemented investing (domain logic + command logic only).

The APIs already available are the following, and uses JSON as their body:

- `POST /` creates a new loan. Since this service focuses on loans only, I'd assume it'd be deployed under some other /loans or http://loan name.
- `POST /approve/:loan_id` approves a loan.
- `POST /invest/:loan_id` invests on a loan. 

This could be more RESTful by treating "loans" as purely an HTTP resource and modifying the state with {"state":"approved"}. However, I thought that it offers little benefit compared to path-level clarity, which would show up in logs and would be easier to maintain.

The project is structured using a variation of clean architecture. Notable points are:

- Sociable unit tests that lets the SUT interacts with its dependencies, to assert the SUT's behavior (e.g. propose command modifies DB)
- Domain logic isolation -- all domain logic is isolated from infrastructure.
- Dependency Inversion principle -- application/business logic accepts interface as its dependency and lets infrastructure layer implement it.

