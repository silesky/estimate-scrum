
Simple flow
1. user goes to www.foo.com and creates game and gets a sessionID.
2. users open www.foo.com/sessionID
3. all users can:
  - set available points
  - change issue
  - vote

... basically, all users are admins. Why? Well, what if a user drops out?


Complex Flow:
1. admin goes to www.foo.com and creates game and gets a sessionID and authGUID.
2. admin clicks 'new game' and is redirected to admin page of www.foo.com/sessionID?auth=authGUID.
3. users open www.foo.com/sessionID

User Privileges:
- vote
- enter name

Admin Privileges:
- vote
- enter name
- change issue
- set available points



