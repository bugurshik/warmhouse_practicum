// Sensor represents a smart home sensor
namespace TelemetryAPI;

public record SensorUpdated
{
    public int Id { get; set; }
    public string Type { get; set; }
    public string Value { get; set; }
    public string Unit { get; set; }
    public string Status { get; set; }

    public TimeSpan TimeSpan { get; set; }
}