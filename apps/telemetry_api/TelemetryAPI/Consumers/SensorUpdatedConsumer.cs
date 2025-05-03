using MassTransit;

namespace TelemetryAPI.Consumers;

public class SensorUpdatedConsumer(ITelemetryStorage storage) : IConsumer<SensorUpdated>
{
    [EndpointName("temperature-t")]
    public Task Consume(ConsumeContext<SensorUpdated> context)
    {
        storage.Add(context.Message);
        Console.WriteLine($"{context.Message.ToString() } Taked!");
        return Task.CompletedTask;
    }
}