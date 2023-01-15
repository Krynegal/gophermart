# Gophermart
The system is an HTTP API with the following business logic requirements:
- registration, authentication and authorization of users;
- receiving order numbers from registered users;
- accounting and maintaining a list of transferred order numbers of a registered user;
- accounting and maintenance of a registered user's savings account;
- verification of accepted order numbers through the loyalty points calculation system;
- accrual of the required remuneration for each suitable order number to the user's loyalty account.
    
 # endpoints
 The Hoffermart cumulative loyalty system provides the following HTTP handlers:

 POST /api/user/register — user registration;
 POST /api/user/login — user authentication;
 POST /api/user/orders — user uploading the order number for calculation;
 GET /api/user/orders — getting a list of order numbers uploaded by the user, their processing statuses and information about charges;
 GET /api/user/balance — getting the current account balance of the user's loyalty points;
 POST /api/user/balance/withdraw — request to deduct points from the savings account to pay for a new order;
 GET /api/user/balance/withdrawals — getting information about the withdrawal of funds from the savings account by the user.
 
 # environment
 The service supports configuration by the following methods:
- address and port of service launch: OS environment variable RUN_ADDRESS or flag -a;
- database connection address: the OS environment variable DATABASE_URI or the -d flag;
- the address of the accrual calculation system: the ACCRUAL_SYSTEM_ADDRESS OS environment variable or the -r flag.
