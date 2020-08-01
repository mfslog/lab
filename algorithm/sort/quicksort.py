#!usr/bin/env python3
#coding=utf8



source = [1,5,2,8,10,3,7,6,9,4]


class Solution:
    def sortArray(self, arry:list)->list:
        length = len(arry)
        if length == 1 or length == 0:
            return arry

        mid = int(length/2)
        midNum = arry[mid]
        #print("mid:%d,midNum:%d"%(mid,midNum))
        left = []
        right = []
        #print(arry[:mid])
        #print(arry[mid+1:])
        for num in arry[:mid]:
            if num >= midNum:
                right.append(num)
            else:
                left.append(num)

        for num in arry[mid+1:]:
            if num >= midNum:
                right.append(num)
            else:
                left.append(num)

        #print("left:",left)
        #print("right:",right)
        if len(left) > 0:
            left = self.sortArray(left)
        if len(right)> 0:
            right = self.sortArray(right)
        left.append(midNum)
        return left + right

def main():
    s = Solution()
    result =s.sortArray(source)
    print(result)

if __name__ == "__main__":
    main()