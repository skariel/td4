<script>
    export let test;
</script>

<style>
    .main {
        display:       grid;
        border:        1px solid #333333;
        border-radius: 15px;
        padding:       10px;
        cursor:        pointer;

        grid-template-rows: 30px 25px auto 30px 50px;        
    }

    .main:hover {
        border:        3px solid rgb(255,62,0);
        padding: 8px;
    }

    .avatar {
        width:    40px;
        height:   40px;
        margin-right:20px;        
    }

    .displayname {

    }
    .topbar {
        grid-column: 1 / 5;
        grid-row: 1;
        display: flex;
    }

    .testid {
        margin-left: auto;
    }

    .testtitle {
        grid-column: 1 / 5;
        grid-row: 2;
    }

    .testdescr {
        grid-column: 1 / 5;
        grid-row: 3;
    }

    .updated {
        grid-column: 1 / 5;
        grid-row: 4;
    }

    .failed {
        grid-column: 1;
        grid-row: 5;
        padding: 5px;
    }

    .passed {
        grid-column: 2;
        grid-row: 5;
        padding: 5px;
    }

    .pending {
        grid-column: 3;
        grid-row: 5;
        padding: 5px;
    }

    .wip {
        grid-column: 4;
        grid-row: 5;
        padding: 5px;
    }
</style>

<div class="main" onclick="location.href='{"/test/"+test.id+"/0"}';">

    <div class="topbar">
        <img class="avatar" src={test.avatar} alt="avatar"/>
        <p class="displayname">{test.display_name}</p>
        <p class="testid">{test.id}</p>
    </div>
    <p class="testtitle">Title: {test.title}</p>
    <p class="testdescr">Descr: {test.descr}</p>
    <p class="updated">u: {test.ts_updated}</p>

    <p class="failed">fail {test.total_fail}</p>
    <p class="passed">pass {test.total_pass}</p>
    <p class="pending">pending {test.total_pending}</p>
    <p class="wip">wip {test.total_wip}</p>

</div>


<!-- 

avatar: "https://avatars2.githubusercontent.com/u/2347816?v=4"
​
code: "lkj"
​
created_by: "github:2347816"
​
descr: "lkj"
​
display_name: "skariel"
​
id: 9
​
title: "lkj"
​
total_fail: 0
​
total_pass: 0
​
total_pending: 0
​
total_wip: 0
​
ts_created: "2020-03-06T00:32:59.065238+02:00"
​
ts_updated: "2020-03-06T00:32:59.065238+02:00"
​
updated_by: "github:2347816" -->
