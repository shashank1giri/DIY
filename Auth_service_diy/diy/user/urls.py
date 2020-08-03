from django.urls import path,include
from django.contrib import auth
from .views import *
from django.contrib.auth import views as auth_views
from django.conf.urls import url

app_name = "user"
urlpatterns = [
    path('',home,name = 'home'),
    url('login/',auth_views.LoginView.as_view(template_name = 'login.html'),name = 'login'),
    url('logout/',auth_views.LogoutView.as_view(),name = 'logout'),
    url('register/',register,name = 'register'),
    url('password_change/',password_change, name='set_password'),
    url('all/',display_all_users,name = 'all'),
    url('(?P<id>[0-9]+)/',display_user_details, name = 'show'),

]