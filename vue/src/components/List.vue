<template>
    <div>
        <div :class="[todo.isComplete ? 'complete':'']" id="todo-container" :key="todo.id" v-for="(todo,index) in todos" >
            
            <h3 >{{ todo.title }}</h3>
             <input @click="$emit('update-item',index)" id="check" type="checkbox"/>
            <p>{{ 
               todo.timestamp
                
            }} </p>
            
            <button @click="$emit('delete-item',index)">x</button>
           
        </div>
    </div>
</template>

<script>
export default {
    name:'List',
    props:{
        todos:Array
    },
}
</script>

<style scoped>
    #todo-container{
        background: rgb(212, 212, 212);
        padding: 10px;
        margin:5px 0;
         border-left: solid 4px transparent;
    }

    h3{
        font-size: 0.9em;
        width: 90%;
        display: inline-block;
    }
    p{
        font-size: 0.8em;;
    }

    .complete {
       border-left: solid 4px rgb(0, 204, 0) !important;

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

    #check{
        padding:10px;
        height: 15px;
        width: 15px;
    }
    
</style>