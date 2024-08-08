import "github.com/google/generative-ai-go/genai"
import "google.golang.org/api/option"

ctx := context.Background()
client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
model := client.GenerativeModel("gemini-1.5-flash")

resp, err := model.GenerateContent(
ctx,
genai.Text("What's is in this photo?"),
genai.ImageData("jpeg", imgData))
