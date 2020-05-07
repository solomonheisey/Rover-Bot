import os
import sys
import time

from TwitterAPI import TwitterAPI

sys.stdout.flush()
time.sleep(1)

message = sys.argv[1]
oauth_token = os.environ['ACCESS_TOKEN']
oauth_token_secret = os.environ['ACCESS_TOKEN_SECRET']
api_key = os.environ['CONSUMER_KEY']
api_secret = os.environ['CONSUMER_SECRET']

api = TwitterAPI(api_key, api_secret, oauth_token, oauth_token_secret)

file = open('Mars.jpg', 'rb')
data = file.read()
r = api.request('statuses/update_with_media', {'status': message}, {'media[]':data})
print ("Status Updated!")


