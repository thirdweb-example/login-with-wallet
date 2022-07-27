from flask import Flask, request, make_response, jsonify
from thirdweb.types import LoginPayload
from thirdweb import ThirdwebSDK
from datetime import datetime, timedelta
import os

app = Flask(__name__)

@app.route("/login", methods=["POST"])
def login():
    private_key = os.environ.get("ADMIN_PRIVATE_KEY")
    if not private_key:
        print("Missing ADMIN_PRIVATE_KEY environment variable")
        return "Admin private key not set", 400

    sdk = ThirdwebSDK.from_private_key(private_key, "mumbai")
    payload = LoginPayload.from_json(request.json["payload"])

    # Generate an access token with the SDK using the signed payload
    domain = "thirdweb.com"
    token = sdk.auth.generate_auth_token(domain, payload)

    res = make_response()
    res.set_cookie(
        "access_token", 
        token,
        path="/",
        httponly=True,
        secure=True,
        samesite="strict"
    )
    return res, 200

@app.route("/authenticate", methods=["POST"])
def authenticate():
    private_key = os.environ.get("ADMIN_PRIVATE_KEY")
    if not private_key:
        print("Missing ADMIN_PRIVATE_KEY environment variable")
        return "Admin private key not set", 400

    sdk = ThirdwebSDK.from_private_key(private_key, "mumbai")

    # Get access token off cookies
    token = request.cookies.get("access_token")
    if not token:
        return "Unauthorized", 401
    
    domain = "thirdweb.com"

    try:
        address = sdk.auth.authenticate(domain, token)
    except:
        return "Unauthorized", 401
    
    return jsonify(address), 200

@app.route("/logout", methods=["POST"])
def logout():
    res = make_response()
    res.set_cookie(
        "access_token", 
        "none",
        expires=datetime.utcnow() + timedelta(seconds=5)
    )
    return res, 200