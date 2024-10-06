# To_Do_List
## Description
The To_Do_List is a simple and interactive web application designed to help users manage their tasks efficiently. Built with HTML, CSS, and JavaScript, this application demonstrates CRUD (Create, Read, Update, Delete) operations through a user-friendly interface. Users can fetch, update, and delete tasks based on their name and selected date, making task management straightforward and organized.
## Features
- Add Task: Users can add a new task by entering the task name, date, and description.
* Fetch Tasks: Retrieve and view existing tasks based on specific criteria.
+ Update Task: Modify the details of an existing task.
- Delete Task: Remove tasks from the list as needed.
## Installation Instructions
To set up the To_Do_List application locally, follow these steps:
- Install PostgreSQL: Download and install the latest version of PostgreSQL from the official website.
* Install Go (Golang): Download and install Golang from the official site.
+ Choose a Code Editor: Use any code editing tool of your choice (e.g., Visual Studio Code, Sublime Text, etc.).
- Set Up Environment Variables: Create a .env file to store your database credentials and other necessary environment variables.
## Usage
Once the application is set up, you can start using it:
- Run the Golang backend server.
* Open your web browser and navigate to the application. You will see four buttons: Add, Fetch, Update, and Delete.
+ ***Add***: Click to add a new task by providing the necessary details like name, date, and task description.
- ***Fetch***: Use this button to view tasks that match specific criteria.
* ***Update***: Modify the details of a task by clicking this button.
+ ***Delete***: Remove a task by selecting it and clicking the delete button.
- You can cross-check the tasks in your PostgreSQL database to confirm the CRUD operations were successful.
## Technologies Used
The application is built using the following technologies:
### Frontend:
- HTML
* CSS
+ JavaScript
### Backend:
- Golang
* PostgreSQL for database management
### Libraries:
- database/sql: Implemented robust database interactions.
* encoding/json: Developed API integrations handling JSON data.
+ log: Enhanced application reliability through structured logging.
- net/http: Built RESTful services.
* os: Managed configuration and file operations securely.
+ strings: Executed complex string manipulations for data processing.
- godotenv: Utilized for environment variable management in development.
## Getting Started
- Clone the repository:
```git clone <your-repo-url>```

```cd To_Do_List```

* Install required Go packages: ```go mod tidy```

+ Set up your PostgreSQL database and update your .env file with the appropriate database credentials.

- Run the application: ```go run main.go```

* Open your web browser and go to (http://localhost:8080) to access the application.