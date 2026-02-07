# gRPC Blog Service – API Documentation

## Service: BlogService

### CreatePost
**Input**
- title (string)
- content (string)
- author (string)
- publication_date (timestamp)
- tags ([]string)

**Output**
- BlogPost on success
- error string on failure

### ReadPost
**Input**
- post_id (string)

**Output**
- BlogPost if found
- error string if not found

### UpdatePost
**Input**
- post_id (string)
- title (string)
- content (string)
- author (string)
- tags ([]string)

**Output**
- Updated BlogPost
- error string on failure

### DeletePost
**Input**
- post_id (string)

**Output**
- success (bool)
- error string if deletion fails
