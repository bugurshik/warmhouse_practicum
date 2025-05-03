namespace TelemetryAPI;

public class MyBackgroundService(HttpClient client) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (true)
        {
            //Триггер обновления сенсоров
            try
            {
                var response = await client.GetAsync("http://smarthome-app:8080/api/v1/sensors", stoppingToken);
                Console.WriteLine(response.StatusCode.ToString());
            }
            catch (Exception ex)
            {
                Console.WriteLine(ex.Message);
            }

            await Task.Delay(30000, stoppingToken);
        }
    }
}