using MassTransit;
using Microsoft.AspNetCore.Mvc;
using TelemetryAPI;
using TelemetryAPI.Consumers;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddHostedService<MyBackgroundService>();
builder.Services.AddHttpClient();
builder.Services.AddSingleton<ITelemetryStorage, FakeTelemetryStorage>();
builder.Services.AddMassTransit(config =>
{
    config.AddConsumer<SensorUpdatedConsumer>();
    config.UsingRabbitMq((context, cfg) =>
    {
        cfg.Host("my-rabbit", "/", h =>
        {
            h.Username("guest");
            h.Password("guest");
        });

        cfg.ConfigureEndpoints(context);
        cfg.ReceiveEndpoint("temperature-t", e =>
        {
            e.UseRawJsonSerializer();
            e.ConfigureConsumer<SensorUpdatedConsumer>(context);
        });
    });
});

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.MapGet("history", 
    ([FromServices] ITelemetryStorage storage,
    [FromQuery] string? location = null,
    [FromQuery] DateTime? fromDate = null,
    [FromQuery] DateTime? toDate = null,
    [FromQuery] string? sensorType = null,
    [FromQuery] string? sensorId = null) =>
{
    var sensors = storage.GetAll();
    if (sensorId is not null)
        sensors = sensors.Where(x => x.SensorId == sensorId);

    return TypedResults.Ok(sensors);
});

app.UseRouting();

app.Run();