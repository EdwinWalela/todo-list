<template>
    <div>
        <div id="todo-container" :key="todo.id" v-for="todo in todos" :class="todo.isComplete">
            
            <h3 :class="todo.isComplete">{{ todo.title }}</h3>
            <p>{{ 
                new Date(todo.timestamp *1000).toString().slice(0,16) 
                
            }} </p>
            
            <button>x</button>
            <button>Complete</button>
        </div>
    </div>
</template>

<script>
export default {
    name:'List',
    props:{
        todos:Array
    }
}
</script>

<style scoped>
    #todo-container{
        background: gainsboro;
        padding: 10px;
        margin:5px 0;
    }

    h3{
        font-size: 0.9em;
    }
    p{
        font-size: 0.8em;;
    }

    .true {
        text-decoration: line-through;

    }

    button:nth-of-type(1){
        border:none;
        background: red;
        color: white;
        display: block;
        margin: 1em 0;
    }
     button:nth-of-type(2){
        border:none;
        background: green;
        color: white;
    }
    
</style>