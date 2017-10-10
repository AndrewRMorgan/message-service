# Message Storage Microservice

## User stories:
1. I can store a message as a string and be given an id for that message. 
2. I can retrieve a previously saved string using it's unique id.
---
### Example Input:
To store a message using a curl command:
```text
curl https://fierce-shelf-71001.herokuapp.com/messages/ -d "Your message here"
```
### Example Output:
```js
{ "id": 12345 }
```
---
### Example Input:
To retrieve a previously saved message using a curl command:
```text
curl https://fierce-shelf-71001.herokuapp.com/messages/12345
```
### Example Output:
```js
Your message here
```
