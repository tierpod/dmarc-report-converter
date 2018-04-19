$(document).ready(function(){
    console.log('search data and labels inside table');

    // var max_items = 10;
    var data = [];
    var labels = [];

    $('#items-table tbody tr').each(function() {
        if (max_items == 0) {
            return false;
        };

        var hostname = $(this).find('#hostname').text();
        if (hostname == "") {
            hostname = $(this).find('#ip').text();
        };
        var msgc = $(this).find('#msgc').text();
        labels.push(hostname);
        data.push(msgc);

        //max_items--;
    });

    var ctx = document.getElementById('stats-chart').getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                label: 'hostnames',
                data: data,
                backgroundColor: [
                    'rgba(255, 99, 132, 0.4)',
                    'rgba(54, 162, 235, 0.4)',
                    'rgba(255, 206, 86, 0.4)',
                    'rgba(75, 192, 192, 0.4)',
                    'rgba(153, 102, 255, 0.4)',
                    'rgba(255, 159, 64, 0.4)'
                ],
                borderColor: 'rgba(0, 0, 0, 0.1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:true
                    }
                }]
            }
        }
    });
});
