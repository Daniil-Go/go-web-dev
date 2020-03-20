<footer>
   <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.1.5/js/uikit.min.js"></script>
   <script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.1.5/js/uikit-icons.min.js"></script>

      <script>
           function check(checkboxElem) {
               var itemID = checkboxElem.id;

               if (checkboxElem.checked) {
                   alert("Закрыта задача " + itemID)
               } else {
                   alert("Открыта задача " + itemID)
               }
           }
           function newItem() {
               alert("Событие создания новой задачи")
           }
       </script>

</footer>