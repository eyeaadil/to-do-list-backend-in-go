package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"to-do-list-api/config"
	"to-do-list-api/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo"
)

// var taskCollection *mongo.Collection = config.DB.Collection("tasks")
var taskCollection *mongo.Collection

// GetTasks retrievs all tasks from the database

// GetTasks function HTTP request ka response handle karne ke liye banayi gayi hai
func GetTasks(w http.ResponseWriter, r *http.Request) {
	taskCollection = config.DB.Collection("tasks")

	// Pehle ek empty tasks slice banate hain, jo models.Task type ka hoga
	var tasks []models.Task

	// taskCollection se saare tasks ko find karne ke liye cursor ka use kar rahe hain
	cursor, err := taskCollection.Find(context.Background(), bson.D{{}})

	// Agar koi error aati hai toh yeh handle karega, aur client ko 500 Internal Server Error return karega
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}

	// Jab tak kaam ho raha hai, tab tak cursor ko background mein close nahi karte,
	// lekin function ke khatam hote hi cursor ko close kar dete hain
	defer cursor.Close(context.Background())

	// Loop chalayenge jab tak cursor ke next document ko read nahi kar lete
	for cursor.Next(context.Background()) {
		var task models.Task

		// Jo current task document read ho raha hai, usko decode karte hain
		// Agar decode mein koi error aayi, toh error handle karenge aur return kar denge
		if err := cursor.Decode(&task); err != nil {
			http.Error(w, "Failed to decode task", http.StatusInternalServerError)
			return
		}

		// Successfully decode hone par task ko tasks slice mein append karenge
		tasks = append(tasks, task)
	}

	// Response ke headers mein content type ko "application/json" set karenge
	w.Header().Set("Content-Type", "application/json")

	// Finally, tasks ko JSON format mein client ko send karenge
	json.NewEncoder(w).Encode(tasks)
}

// CreateTask creates a new task in the database
func CreateTask(w http.ResponseWriter, r *http.Request) {

	taskCollection = config.DB.Collection("tasks")
    // Pehle 'task' naam ka ek variable declare karte hain jo models.Task type ka hoga
    // Yeh task data ko request body se rakhne ke liye use kar rahe hain
    var task models.Task

    // JSON request body ko 'task' variable mein decode karte hain
    // Agar decoding fail ho gayi (e.g., JSON format galat hai), toh 400 Bad Request error return karenge
    if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
        // Agar request body invalid hai toh 400 Bad Request response bhejte hain
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return // Decoding fail hone ke baad function ko exit kar dete hain
    }

    // Ek unique ObjectID assign karte hain task ko, jo task ka identifier hoga database mein
    task.ID = primitive.NewObjectID()

    // MongoDB collection mein task ko insert karte hain
    // 'context.Background()' request ke lifetime ko manage karne ke liye use hota hai
    // Agar insertion mein koi error aayi, toh usko 'err' variable mein capture karenge
    _, err := taskCollection.InsertOne(context.Background(), task)

    // Agar insertion ke dauran koi error aayi (e.g., database issue), toh 500 Internal Server Error return karenge
    if err != nil {
        // Agar task insert nahi ho paya, toh 500 Internal Server Error response bhejte hain
        http.Error(w, "Failed to create task", http.StatusInternalServerError)
        return // Task creation fail hone ke baad function ko exit kar dete hain
    }

    // Agar sab kuch successfully ho gaya, toh 201 Created status code return karenge
    w.WriteHeader(http.StatusCreated)

    // Created task ko JSON format mein encode karke response body mein bhejte hain
    json.NewEncoder(w).Encode(task)
}

// DeleteTask deletes a task by ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskCollection = config.DB.Collection("tasks")
    // Extract the URL parameters (specifically, the task ID) from the request using mux.Vars.
    params := mux.Vars(r)

    // Convert the task ID from a string (in the URL) to a MongoDB ObjectID.
    // This is necessary because MongoDB uses ObjectIDs to identify documents.
    id, err := primitive.ObjectIDFromHex(params["id"])

    // Check if the task ID is invalid (e.g., not a valid ObjectID format).
    if err != nil {
        // If the ID is invalid, return a 400 Bad Request HTTP response with an error message.
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return // Exit the function early because the ID is invalid.
    }

    // Create a filter that specifies which document (task) to delete based on its "_id".
    filter := bson.M{"_id": id}

    // Attempt to delete the task that matches the filter (i.e., the task with the given ObjectID).
    result, err := taskCollection.DeleteOne(context.Background(), filter)

    // Check if there was an error during deletion or if no document was deleted (DeletedCount == 0).
    if err != nil || result.DeletedCount == 0 {
        // If the deletion failed or no task was found to delete, return a 500 Internal Server Error.
        http.Error(w, "Failed to delete task", http.StatusInternalServerError)
        return // Exit the function early due to the error.
    }

    // If everything is successful, return a 204 No Content status code, indicating the task was deleted.
    w.WriteHeader(http.StatusNoContent)
}

// UpdateTask ek existing task ko ID ke basis pe update karta hai
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskCollection = config.DB.Collection("tasks")
    // URL parameters se task ID ko extract kar rahe hain, mux.Vars ka use karke.
    params := mux.Vars(r)

    // Task ID (jo string hai) ko MongoDB ObjectID mein convert kar rahe hain.
    // Agar task ID URL mein invalid ho, toh 400 Bad Request ka response bhejte hain.
    id, err := primitive.ObjectIDFromHex(params["id"])
    if err != nil {
        // Agar ID conversion fail ho gaya, toh error denge aur function ko yahi exit karenge.
        http.Error(w, "Invalid task ID", http.StatusBadRequest)
        return
    }

    // Ek variable declare kar rahe hain jo updated task details ko request body se hold karega.
    var updatedTask models.Task

    // JSON request body ko updatedTask struct mein decode kar rahe hain.
    // Agar decoding fail ho jaaye (jaise ki agar JSON valid na ho), toh 400 Bad Request denge.
    if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
        // Agar request payload mein error ho, toh function yahi exit karenge.
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Ek filter banate hain jo task ko uski ObjectID ke basis pe find karega (update karne ke liye).
    filter := bson.M{"_id": id}

    // Update object bana rahe hain jo $set operator ka use karta hai, jisse task ke fields update hote hain.
    // Hum "name", "description", aur "status" fields ko update kar rahe hain, jo request mein aaye hain.
    update := bson.M{
        "$set": bson.M{
            "name":        updatedTask.Name,
            "description": updatedTask.Description,
            "Completed" : updatedTask.Completed,
            "status":      updatedTask.Status,
        },
    }

    // Task document par update operation perform kar rahe hain jo filter (task ID) se match karega.
    // Agar koi task filter se match nahi karta ya koi error aata hai, toh 500 Internal Server Error denge.
    result, err := taskCollection.UpdateOne(context.Background(), filter, update)
    if err != nil || result.MatchedCount == 0 {
        // Agar update fail hota hai ya koi task nahi milta, toh function exit karenge.
        http.Error(w, "Failed to update task", http.StatusInternalServerError)
        return
    }

    // Response header set karte hain jisse content type JSON ho.
    w.Header().Set("Content-Type", "application/json")

    // Updated task ko JSON format mein encode kar ke response mein bhejte hain.
    json.NewEncoder(w).Encode(updatedTask)
}
