// Sensor represents a smart home sensor
using System.Text.Json.Serialization;
using Newtonsoft.Json;

namespace TelemetryAPI.Consumers;

public record SensorUpdated
{
    [JsonProperty("sensor_id")]
    [JsonPropertyName("sensor_id")]
    public string SensorId { get; set; }

    [JsonProperty("sensor_type")]
    [JsonPropertyName("sensor_type")]
    public string SensorType { get; set; }
    public double Value { get; set; }
    public string Unit { get; set; }
    public string Status { get; set; }

    public DateTime Timestamp { get; set; }
}