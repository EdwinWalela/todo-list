<template>
    <div id="task-container">
        <h1>Todo Manager</h1>
        <form id="form" @submit="onSubmit">
            <label for="">Task</label>
            <input type="task" v-model="form.title" placeholder="Task title"/>
            <br>
            <label for="">Day & Time</label>
            <input  type="date" v-model="form.date" placeholder="Due date">
            <br>
            <input id="submit-btn" type="submit" value="Add Task">
        </form>

        <List 
            :todos="todos"
            @delete-item="deleteItem"
            @update-item="updateItem"
        />
        
    </div>
</template>

<script>

import List from './List.vue'

export default {
    name:"Todo",
    components:{
        List
    },
    created(){
        // Fetch 
        this.todos = [
               
            ]
    },

    methods: {
        async onSubmit(e){
            e.preventDefault();
            if(this.form.title == ""){
                alert("Task title required")
                return
            }
            if(typeof this.form.date == "undefined"){
                 alert("Date required")
                return
            }

            this.todos.push({
                title:this.form.title,
                isComplete:false,
                timestamp:this.form.date,
                })
           
           this.form.title = ""
        },

        deleteItem(index){
            this.todos.splice(index,1);
        },

        updateItem(index){
            this.todos[index].isComplete  = ! this.todos[index].isComplete
        }
    },
    data(){
    return {
      todos:[
         
      ],
      form:{
          title:'',
          timestamp:'',
          isComplete:false,
      }
     
    }
  },

}
</script>

<style scoped>
    *{
        
        font-family: sans-serif
    }
    h1{
        text-align: center;
        font-size: 1.4em;
    }

    #task-container{
         box-shadow: 0px 1px 3px rgba(0, 60, 255, 0.5);
        width: 50%;
        max-width: 400px;
        margin: 2em auto;
        padding: 1em 2em;
    }

    label{
        font-size: 0.8em;
        font-weight: 600;
    }

    input{
        width:100%;
        padding: 5px;
        font-size: 0.9em;
        margin: 0.8em 0;
        letter-spacing: 1.3px;
    }

    #submit-btn{
        background: cornflowerblue;
        color:white;
        border:none;
        padding: 8px;
    }

</style>