using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.OpenApi.Models;
using Authenticator.Services;
using Authenticator.Database.Context;
using Authenticator.Interfaces;
using Authenticator.Models.Auth;
using Authenticator.Models.User;
using System;
using System.IO;
using System.Reflection;
using Authenticator.Models;
using Authenticator.Middleware;

namespace Authenticator
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        // This method gets called by the runtime. Use this method to add services to the container.
        public void ConfigureServices(IServiceCollection services)
        {
            services.AddCors(options =>
            {
                options.AddPolicy("AllowAllOrigins",
                    builder => builder.AllowAnyOrigin());
            });


            //for postgresql
            var user   = Environment.GetEnvironmentVariable("PG_USER");
            var pass   = Environment.GetEnvironmentVariable("PG_PASS");
            var server = Environment.GetEnvironmentVariable("PG_SERVER");
            var db     = $"{File.ReadAllText("/run/secrets/kubernetes.io/serviceaccount/namespace")}.authenticator";

            services.AddDbContext<GlobalDbContext>(options =>
                options.UseNpgsql($"Host={server};Port=5432;Database={db};Username={user};Password={pass}"),
                ServiceLifetime.Transient
            );
            services.AddDatabaseDeveloperPageExceptionFilter();

            services.AddIdentity<UserAccount, IdentityRole<int>>()
                .AddDefaultTokenProviders()
                .AddEntityFrameworkStores<GlobalDbContext>();


            services.Configure<IdentityOptions> (options => {
                options.Password.RequireDigit = true; 
                options.Password.RequireLowercase = true; 
                options.Password.RequireNonAlphanumeric = true; 
                options.Password.RequireUppercase = true; 
                options.Password.RequiredLength = 8; 
                options.Password.RequiredUniqueChars = 1; 

                options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5); 
                options.Lockout.MaxFailedAccessAttempts = 5; 
                options.Lockout.AllowedForNewUsers = true;

                options.User.AllowedUserNameCharacters = 
                    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
                options.User.RequireUniqueEmail = true;  

                options.SignIn.RequireConfirmedEmail = false;            
                options.SignIn.RequireConfirmedPhoneNumber = false;     
            });

            services.AddControllers();
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo
                {
                    Title = "Host",
                    Version =
                    "v1"
                });
                c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme()
                {
                    Name = "Authorization",
                    Type = SecuritySchemeType.ApiKey,
                    Scheme = "Bearer",
                    BearerFormat = "JWT",
                    In = ParameterLocation.Header,
                    Description = "JWT Authorization header using the Bearer scheme.",
                });

                c.AddSecurityRequirement(new OpenApiSecurityRequirement
                {
                    {
                        new OpenApiSecurityScheme
                            {
                                Reference = new OpenApiReference
                                {
                                    Type = ReferenceType.SecurityScheme,
                                    Id = "Bearer"
                                }
                            },
                            new string[] {}
                    }
                }); 
            });

            services.Configure<JwtOptions>(Configuration.GetSection("JwtOptions"));
            services.Configure<SystemConfig>(Configuration.GetSection("SystemConfig"));
            services.Configure<IdentityOptions>(options =>
            {
                options.Password.RequireDigit = false;
                options.Password.RequiredLength = 5;
                options.Password.RequireLowercase = true;
                options.Password.RequireNonAlphanumeric = false;
                options.Password.RequireUppercase = false;
            });

            services.AddTransient<ITokenGenerator, TokenGenerator>();
            services.AddMvc();
        }

        // This method gets called by the runtime. Use this method to configure the HTTP request pipeline.
        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            app.UseSwagger();
            app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "signalling v1"));

            app.UseCors(x => x
                .AllowAnyMethod()
                .AllowAnyHeader()
                .WithMethods("GET", "POST")
                .AllowCredentials()
                .SetIsOriginAllowed(origin => true)); // allow any origin


            app.UseRouting();
            app.UseMiddleware<JwtMiddleware>();
            app.UseMiddleware<AuthorizeMiddleWare>();
            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllerRoute(
                    name: "default",
                    pattern: "{controller=Home}/{action=Index}/{id?}");
            });
        }
    }
    
}