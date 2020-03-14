<script>
    export let test;
</script>

<style>
    .main {
        display:       grid;
        background-color: rgb(236,230,226);
        border:        1px solid #ffffff;
        border-radius: 7px;
        padding:       10px;
        /* cursor:        pointer; */

        grid-template-columns: 25% 25% 25% 25%;        
        grid-template-rows: auto auto auto auto auto;        
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
    }

    .passed {
        grid-column: 2;
        grid-row: 5;
    }

    .pending {
        grid-column: 3;
        grid-row: 5;
    }

    .wip {
        grid-column: 4;
        grid-row: 5;
    }
</style>

<!-- <div class="main" onclick="location.href='{"/test?id="+test.id+"&page=0"}';"> -->
<div class="main">

    <div class="topbar">
        <img class="avatar" src={test.avatar} alt="avatar"/>
        <h4 class="displayname">{test.display_name}</h4>
        <h4 class="testid">{test.id}</h4>
    </div>
    <h4 class="testtitle">Title: <a href='{"/test?id="+test.id+"&page=0"}'>{test.title}</a></h4>
    <h4 class="testdescr">Descr: {test.descr}</h4>
    <h4 class="updated">u: {test.ts_updated}</h4>

    <h4 class="failed">fail {test.total_fail}</h4>
    <h4 class="passed">pass {test.total_pass}</h4>
    <h4 class="pending">pending {test.total_pending}</h4>
    <h4 class="wip">wip {test.total_wip}</h4>

</div>

