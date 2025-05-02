using MassTransit;
using TelemetryAPI;

public class TelemetryConsumer : IConsumer<SensorUpdated>
{
    public async Task Consume(ConsumeContext<SensorUpdated> context)
    {
        var message = context.Message;
        Console.WriteLine($"User created: {message.Id}");
    }
}