# Go Server Setup

This guide will walk you through setting up the Go server with Redis and MySQL.

## 1. Configure the Redis Queue

To handle background tasks, we've using Redis queues. Follow these steps to configure Redis:

### Step 1: Create Redis Channels
You need to create two channels in Redis:

- **ticket_queue**: Handles ticket processing and task reminder emails.
- **gmail_queue**: Handles sending emails from the sales team to customers (Gmail-related tasks).

You can use the Redis CLI or a Redis client to create these channels. For example, run the following commands:

redis-cli
> PUBLISH ticket_queue "task_data"
> PUBLISH gmail_queue "task_data"

## 2. Configure MYSQL
We have to configure the mysql credentials in the .env file just like below

DB_USER = ""
DB_PASSWORD = ""
DB_HOST = ""
DB_NAME = ""

Here you can add your mysql credentials in the env file, And import the .sql file in your Database.


## 3. Gmail API Integration

To send emails and implement login/signup via Google, you will need to set up OAuth 2.0 credentials from Google Cloud and enable the Gmail API. Follow the steps below:

### Step 1 :
    Go to Google Cloud Console:
        Visit the Google Cloud Console.

    Create a New Project:

        Click the project dropdown in the top bar and click "New Project."

        Give your project a name and click "Create."

### Step 2 :
    Enable Gmail API
        Once your project is created, go to the API Library in the Google Cloud Console.

        Search for “Gmail API” and click on it.

        Click "Enable" to enable the Gmail API for your project.

### Step 3: Create OAuth 2.0 Credentials
    Navigate to Credentials:

        1. In the Google Cloud Console, go to the Credentials page.
        Create OAuth 2.0 Client ID:

        2. Click on Create Credentials and select OAuth 2.0 Client ID.
        You may be prompted to configure the OAuth consent screen. Choose External for user type and fill out the necessary fields (app name, email, etc.).
        Configure OAuth 2.0 Client ID:

        3. After configuring the consent screen, click "Create" under the "OAuth 2.0 Client IDs" section.
        For application type, select Web Application.
        Under "Authorized JavaScript origins," add the domain or https://crm.holoware.co for local development.
        Under "Authorized redirect URIs," add the redirect URI where your app will handle OAuth responses (e.g., "Backend URL").

### Step 4: Add the API ID
    Finally add the API key and client secret key in the .env file.
    
## 4. To Start the Serer  
    After all this has been complete you have to run the below command
        go run main.go
