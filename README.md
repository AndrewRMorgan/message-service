# Message Storage Microservice

## User stories:
1. I can store a message as a string and be given an id for that message. 
2. I can retrieve a previously saved string using it's unique id.

## Example Input:

To store a message simply use the curl command:
```text
curl https://fierce-shelf-71001.herokuapp.com/messages/ -d "Your message here"
```

## Example output:
```js
{ "id": 12345 }
```



