from django.shortcuts import render,redirect

# Create your views here.
from django.contrib.auth.decorators import login_required
from django.shortcuts import render
from .forms import RegistrationForm
from django.contrib.auth import authenticate,login,update_session_auth_hash
from django.contrib.auth.forms import PasswordChangeForm
from .models import User
import jwt, datetime, time
from diy.settings import JWT_ALGO, SECRET_KEY


@login_required()
def home(request):
    exp = datetime.datetime.now() + datetime.timedelta(minutes=10)
    exp = time.mktime(exp.timetuple())
    payload = {
        'user_id': request.user.id,
        'exp_time': exp
    }
    print(exp, SECRET_KEY)
    jwt_token = jwt.encode(payload, SECRET_KEY, algorithm=JWT_ALGO)
    print(jwt_token.decode("utf-8"))
    return render(request, 'index.html', {"token":jwt_token.decode("utf-8")})
    #return render(request, 'user/home.html')


def register(request):
    if request.method == "POST":
        form = RegistrationForm(request.POST)
        if form.is_valid():
            user = form.save()
            raw_password = form.cleaned_data.get('password1')
            user = authenticate(request,phone_number = user.phone_number, password = raw_password)
            login(request,user)

        else:
            return render(request,"register.html",{"form":form,"error":"Wrong Entries"})
          
        return redirect('user:home')
    else:
        form = RegistrationForm()
        return render(request,"register.html",{"form":form})


@login_required()     
def password_change(request):
    if request.method == "POST":
        form = PasswordChangeForm(request.user,request.POST)
        if form.is_valid():
            form.save()
            update_session_auth_hash(request,form.user)
            return redirect('user:home')
        else:
            ctx = {"form":form,"error":"Password Not changed"}
            return render(request,"user/password_change.html",ctx)
    else :
        form = PasswordChangeForm(request.user)
        return render(request,"user/password_change.html",{"form":form})


def display_all_users(request):
    users = User.objects.all()
    return render(request,"user/display.html", {"users":users})


def display_user_details(request,id):
    user = User.objects.filter(id = id)
    return render(request,"user/display.html", {"users":user})