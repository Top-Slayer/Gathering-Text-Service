db="./internal/repository/database.db"
table="Categories"

if [[ "$1" == "help" ]]; then
    echo -e "\n[ Categories Database Editor ]"
    echo -e "------------------------------------------------------------"
    echo -e "> $0 <command> [<args>] \n"
    echo -e "Commands: "
    echo -e "  add    <arg> \t add category into database"
    echo -e "  delete <arg> \t delete a category from database"
    echo -e "  show         \t show all categories from database"
    echo -e "------------------------------------------------------------"
    exit 0
elif [[ "$1" == "add" ]]; then
    if [ -z $2 ]; then 
        echo "$0: Missing some argument"
        echo -e "  add    <arg> \t add category into database"
        echo -e "\nEx: $0 add \"body\""
        exit 1
    fi
    if sqlite3 $db "SELECT 1 FROM $table WHERE name='$2' LIMIT 1;" | grep -q 1; then
        echo "$0: The data already included into database";
        exit 1
    fi
    
    sqlite3 $db "INSERT INTO $table(id, name) VALUES((
        SELECT COALESCE(MIN(t1.id + 1), MAX(t2.id) + 1, 1) 
            FROM $table t1 
            LEFT JOIN $table t2 ON t1.id + 1 = t2.id 
            WHERE t2.id IS NULL
        ), '$2')"
    echo "Has inserted data successfully"
    exit 0
elif [[ "$1" == "delete" ]]; then
    if [ -z $2 ]; then 
        echo "$0: Missing some argument"
        echo -e "  delete <arg> \t delete a category from database"
        echo -e "\nEx: $0 delete 1"
        exit 1
    fi
    sqlite3 $db "DELETE FROM $table WHERE id = $2"
    echo "Has deleted id successfully"
    exit 0

elif [[ "$1" == "show" ]]; then
    echo; sqlite3 $db "SELECT * FROM $table" -markdown
    exit 0
else
    echo -e "$0: Try \"help\" argument for guide info"
    exit 1
fi