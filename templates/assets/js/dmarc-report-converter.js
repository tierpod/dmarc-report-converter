$(document).ready(function(){
    console.log('search data and labels inside table');

    var data = [];
    var labels = [];
    var auth_result_pass = 0;
    var auth_result_fail = 0;

    $('#items-table tbody tr').each(function() {
        // count pass and fail auth results
        var msgc = parseInt($(this).find('#msgc').text());

        if ($(this).hasClass('auth-result-pass')) {
            auth_result_pass = auth_result_pass + msgc;
        } else {
            auth_result_fail = auth_result_fail + msgc;
        }

        var hostname = $(this).find('#hostname').text();
        if (hostname == "") {
            hostname = $(this).find('#ip').text();
        };
        var msgc = $(this).find('#msgc').text();
        labels.push(hostname);
        data.push(msgc);
    });

    var ctx = document.getElementById('hosts-chart').getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                label: 'hostnames',
                data: data,
                backgroundColor: [
                    'rgba(255, 99, 132, 0.6)',
                    'rgba(54, 162, 235, 0.6)',
                    'rgba(255, 206, 86, 0.6)',
                    'rgba(75, 192, 192, 0.6)',
                    'rgba(153, 102, 255, 0.6)',
                    'rgba(255, 159, 64, 0.6)'
                ],
                borderWidth: 1
            }]
        },
        options: {
            title: {
                display: true,
                text: 'messages count per hostnames',
            },
            scales: {
                yAxes: [{
                    display: false,
                }]
            },
            legend: {
                display: false,
            },
        }
    });

    var ctx = document.getElementById('stats-chart').getContext('2d');
    var myChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: ['pass', 'fail'],
            datasets: [{
                label: 'stats',
                data: [auth_result_pass, auth_result_fail],
                backgroundColor: [
                    'rgba(44, 160, 44, 0.6)',
                    'rgba(214, 39, 40, 0.6)',
                ],
                borderWidth: 1
            }],
        },
        options: {
            title: {
                display: true,
                text: 'pass/fail messages',
            },
            scales: {
                yAxes: [{
                    display: false,
                }]
            },
            legend: {
                display: false,
                // position: 'right'
            },
        },
    });
});
