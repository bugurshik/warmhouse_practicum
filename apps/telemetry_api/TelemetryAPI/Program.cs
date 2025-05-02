using MassTransit;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddMassTransit(config =>
{
    config.UsingRabbitMq((context, cfg) =>
    {
        cfg.Host("my-rabbit", "/", h =>
        {
            h.Username("guest");
            h.Password("guest");
        });

        // Автоматическая регистрация Consumers (обработчиков сообщений)
        cfg.ConfigureEndpoints(context);
    });
});
builder.Services.AddMemoryCache();

var app = builder.Build();

if (!app.Environment.IsDevelopment())
{
    app.UseExceptionHandler("/Error");
    app.UseHsts();
}

app.UseHttpsRedirection();

app.UseRouting();

app.MapGet("sensors/{sensorId}/history", (string sensorId) => TypedResults.Ok());

app.Run();
