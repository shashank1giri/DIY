from django.db import models
from django.core.validators import RegexValidator
# Create your models here.
from django.contrib.auth.models import AbstractBaseUser,PermissionsMixin,BaseUserManager


class CustomUserManager(BaseUserManager):
    """
    Custom user model manager where phone is the unique identifiers
    for authentication instead of usernames.
    """
    def create_user(self,phone_number, password,**extra_fields):
        """
        Create and save a User with the given phone number and password.
        """
        user = self.model(phone_number= phone_number,**extra_fields)
        user.set_password(password)
        return user

    def create_superuser(self, phone_number, password,**extra_fields):
        """
        Create and save a SuperUser with the given phone number and password.
        """
        user = self.create_user(phone_number,password,**extra_fields)
        user.is_staff = True
        user.is_superuser = True
        user.save()
        return user


class User(AbstractBaseUser,PermissionsMixin):
    name = models.CharField(blank = False,max_length = 20)
    email = models.EmailField(db_index=True)
    phone_regex = RegexValidator(regex=r'^\d{10}$', message="Phone number must be a valid 10 dig number")
    phone_number = models.CharField(max_length = 15,unique = True)
    
    is_superuser = models.BooleanField(default=False)
    is_staff = models.BooleanField(default=False)
    is_active = models.BooleanField(default=True)

    objects = CustomUserManager()
    USERNAME_FIELD = "phone_number"
    REQUIRED_FIELDS = ["name","email"]

    def __str__(self):
        return "%s%s" %(self.email,self.phone_number)
