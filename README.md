# News Service

This is a web-based news management application built with Golang, MongoDB, and HTMx. It provides a clean, interactive UI for creating, reading, updating, and deleting news posts.

---

## Features

- **Create News**: Add a new news post by providing a title and content.
- **Read News**:
  - View all news posts with pagination.
  - Search posts by title or content.
  - View detailed information for a specific news post.
- **Update News**: Edit an existing post's title and content.
- **Delete News**: Remove a post you no longer need.

---

## Routes

### Main Routes

- **`/`**
   - **Method**: `GET`
   - **Description**: Displays a list of all news posts with pagination and search support.
   - **Controller**: `controllers.ReadPosts`

- **`/new`**
   - **Methods**:
     - `GET`: Displays a form for creating a new news post.
     - `POST`: Handles form submission to create a new post.
   - **Controller**: `controllers.CreatePost`

- **`/post/{id}`**
   - **Methods**:
     - `GET`: Displays the details of a single news post.
   - **Controller**: `controllers.ReadPost`

- **`/post/{id}/edit`**
   - **Method**: `POST`
   - **Description**: Handles form submission to update the title or content of an existing post.
   - **Controller**: `controllers.UpdatePost`

- **`/post/{id}/delete`**
   - **Method**: `POST`
   - **Description**: Deletes a specific news post by its ID.
   - **Controller**: `controllers.DeletePost`

- **`/health`**
   - **Method**: `GET`
   - **Description**: Returns a simple "Server is running" response to verify the application is live.

---

## How to Run

### Clone the Repository:
   ```bash
   git clone <repository-url>
   cd <repository-name>
```

### Without Docker 
- Build the application: 
```bash make build ``` 
- Run the application: 
```bash make run ``` 
- Access the application at `http://localhost:8080`. 
### Notes - Ensure MongoDB is running locally and configured to accept connections from the application. - If you encounter any issues, check the logs in your terminal for more details.

### With Docker 
- Build and run the application with Docker Compose: 
```bash make docker-run ``` 
- Access the application at `http://localhost:8080`. 

---

### Adding Sample Data 
- Navigate to the "Create News" page: 
- Go to `http://localhost:8080/new`. 
- Fill in the "Title" and "Content" fields. 
- Submit the form to add the post. 
- Repeat this process to add more posts.

---

### Using the Application 
- **Homepage**: - Open `http://localhost:8080` to view all news posts.
- **Search**: - Use the search bar to filter posts by title or content. 
- **View Single Post**: - Click on a news post title to view its details.
- **Edit Post**: - From a single post's detail page, click "Edit" to modify the post's title or content. 
- **Delete Post**: - From a single post's detail page, click "Delete" to remove the post. 




