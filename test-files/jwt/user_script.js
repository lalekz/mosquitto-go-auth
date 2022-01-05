function checkUser(username, token) {
    if(username == "test" && token == "correct") {
        return true;
    }
    return false;
}

checkUser(username, token);
