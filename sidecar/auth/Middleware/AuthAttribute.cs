using Microsoft.AspNetCore.Mvc.Filters;
using System;

namespace Authenticator.Middleware
{
    [AttributeUsage(AttributeTargets.Method)]
    public class UserAttribute : ActionFilterAttribute
    {
    }

    [AttributeUsage(AttributeTargets.Method)]
    public class AdminAttribute : ActionFilterAttribute
    {
    }
}